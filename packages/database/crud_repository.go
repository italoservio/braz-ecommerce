package database

import (
	"context"
	"fmt"
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
	GetById(collName string, id string) (*mongo.SingleResult, error)
	DeleteById(collName string, id string) (*mongo.DeleteResult, error)
}

func (cr *CrudRepository) GetById(collName string, id string) (*mongo.SingleResult, error) {
	coll := cr.database.Collection(collName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return coll.FindOne(ctx, bson.M{"_id": objectId}), nil
}

func (cr *CrudRepository) DeleteById(collName string, id string) error {
	coll := cr.database.Collection(collName)

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

func (cr *CrudRepository) CreateOne(collName string, structure any) (string, error) {
	coll := cr.database.Collection(collName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := coll.InsertOne(ctx, structure)
	if err != nil {
		return "", err
	}

	fmt.Println(result.InsertedID)

	oid, err := primitive.ObjectIDFromHex(result.InsertedID.(string))
	if err != nil {
		return "", err
	}

	return oid.Hex(), nil
}
