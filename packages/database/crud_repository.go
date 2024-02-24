package database

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CrudRepositoryInterface interface {
	GetById(
		ctx context.Context,
		collection string,
		id string,
		structure any,
	) error
	DeleteById(
		ctx context.Context,
		collection string,
		id string,
	) error
	CreateOne(
		ctx context.Context,
		collection string,
		structure any,
	) (string, error)
	UpdateById(
		ctx context.Context,
		collection string,
		id string,
		structure any,
	) error
	GetPaginated(
		ctx context.Context,
		collection string,
		page int,
		perPage int,
		filters map[string]any,
		projections map[string]int,
		sortings map[string]int,
		structures any,
	) error
	GetByEmail(
		ctx context.Context,
		collection string,
		email string,
		structure any) error
}

type CrudRepository struct {
	logger   logger.LoggerInterface
	database *Database
}

func NewCrudRepository(lg logger.LoggerInterface, db *Database) *CrudRepository {
	return &CrudRepository{logger: lg, database: db}
}

func (cr *CrudRepository) GetById(
	ctx context.Context,
	collection string,
	id string,
	structure any,
) error {
	coll := cr.database.Collection(collection)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	err = coll.FindOne(timeout, bson.M{"_id": objectId}).Decode(structure)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cr.logger.WithCtx(ctx).Error(err.Error())
			return errors.New(exception.CodeNotFound)
		}

		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

func (cr *CrudRepository) DeleteById(
	ctx context.Context,
	collection string,
	id string,
) error {
	coll := cr.database.Collection(collection)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	_, err = coll.UpdateOne(
		timeout,
		bson.M{"_id": objectId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}},
	)

	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

func (cr *CrudRepository) CreateOne(
	ctx context.Context,
	collection string,
	structure any,
) (string, error) {
	coll := cr.database.Collection(collection)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := coll.InsertOne(timeout, structure)
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return "", errors.New(exception.CodeDatabaseFailed)
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (cr *CrudRepository) UpdateById(
	ctx context.Context,
	collection string,
	id string,
	structure any,
) error {
	coll := cr.database.Collection(collection)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	document, err := ParseToDocument(structure)
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	bson := bson.D{{Key: "$set", Value: document}}

	if _, err := coll.UpdateByID(timeout, objectId, bson); err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

type PaginatedOutput interface{}

func (cr *CrudRepository) GetPaginated(
	ctx context.Context,
	collection string,
	page int,
	perPage int,
	filters map[string]any,
	projections map[string]int,
	sortings map[string]int,
	structures any,
) error {
	coll := cr.database.Collection(collection)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filtersBson := mapToBsonM(filters)
	projectionBson := mapToBsonM[int](projections)
	sortingsBson := mapToBsonM[int](sortings)

	limit := int64(perPage)
	skip := int64(perPage * (page - 1))

	cursor, err := coll.Find(timeout, filtersBson, &options.FindOptions{
		Limit:      &limit,
		Skip:       &skip,
		Projection: projectionBson,
		Sort:       sortingsBson,
	})
	if err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, structures); err != nil {
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

func mapToBsonM[T any](m map[string]T) bson.M {
	doc := make(bson.M)
	for k, v := range m {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice {
			doc[k] = bson.M{"$in": v}
		} else {
			doc[k] = v
		}
	}
	return doc
}

func (cr *CrudRepository) GetByEmail(
	ctx context.Context,
	collection string,
	email string,
	structure any,
) error {

	coll := cr.database.Collection(collection)

	cursor, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := coll.FindOne(cursor, bson.M{"email": email}).Decode(structure)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			slog.Error(err.Error())
			return errors.New(exception.CodeNotFound)
		}

		slog.Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}
