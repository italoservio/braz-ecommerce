package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CrudRepository struct{ database *Database }

func NewCrudRepository(db *Database) *CrudRepository {
	return &CrudRepository{database: db}
}

type CrudRepositoryInterface interface {
	GetById(collection string, id string) (*mongo.SingleResult, error)
	DeleteById(collection string, id string) error
	CreateOne(collection string, structure any) (string, error)
	UpdateById(collection string, id string, structure any) error
}

func (cr *CrudRepository) GetById(
	collection string,
	id string,
) (*mongo.SingleResult, error) {
	coll := cr.database.Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return coll.FindOne(ctx, bson.M{"_id": objectId}), nil
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
		return err
	}

	if _, err := coll.DeleteOne(ctx, bson.M{"_id": objectId}); err != nil {
		return err
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
		return "", err
	}

	oid, err := primitive.ObjectIDFromHex(result.InsertedID.(string))
	if err != nil {
		return "", err
	}

	return oid.Hex(), nil
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
		return err
	}

	document, err := ParseToDocument(structure)
	if err != nil {
		return err
	}

	bson := bson.D{{Key: "$set", Value: document}}

	if _, err := coll.UpdateByID(ctx, objectId, bson); err != nil {
		return err
	}

	return nil
}
