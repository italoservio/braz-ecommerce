package database

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CrudRepository struct{ database *Database }

func NewCrudRepository(db *Database) *CrudRepository {
	return &CrudRepository{database: db}
}

type CrudRepositoryInterface interface {
	GetById(collection string, id string, structure any) error
	DeleteById(collection string, id string) error
	CreateOne(collection string, structure any) (string, error)
	UpdateById(collection string, id string, structure any) error
	GetPaginated(collection string, page int, perPage int, filters map[string]any, projections map[string]int, sortings map[string]int, structures any) error
}

func (cr *CrudRepository) GetById(
	collection string,
	id string,
	structure any,
) error {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(structure)
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

func (cr *CrudRepository) DeleteById(
	collection string,
	id string,
) error {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}},
	)

	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

func (cr *CrudRepository) CreateOne(
	collection string,
	structure any,
) (string, error) {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := coll.InsertOne(ctx, structure)
	if err != nil {
		slog.Error(err.Error())
		return "", errors.New(exception.CodeDatabaseFailed)
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (cr *CrudRepository) UpdateById(
	collection string,
	id string,
	structure any,
) error {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	document, err := ParseToDocument(structure)
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	bson := bson.D{{Key: "$set", Value: document}}

	if _, err := coll.UpdateByID(ctx, objectId, bson); err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}

type PaginatedOutput interface{}

func (cr *CrudRepository) GetPaginated(
	collection string,
	page int,
	perPage int,
	filters map[string]any,
	projections map[string]int,
	sortings map[string]int,
	structures any,
) error {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	filtersBson := mapToBsonM(filters)
	projectionBson := mapToBsonM[int](projections)
	sortingsBson := mapToBsonM[int](sortings)

	limit := int64(perPage)
	skip := int64(perPage * (page - 1))

	cursor, err := coll.Find(ctx, filtersBson, &options.FindOptions{
		Limit:      &limit,
		Skip:       &skip,
		Projection: projectionBson,
		Sort:       sortingsBson,
	})
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, structures); err != nil {
		slog.Error(err.Error())
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
