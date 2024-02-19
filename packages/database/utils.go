package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseToDocument(structure any) (*bson.D, error) {
	data, err := bson.Marshal(structure)
	if err != nil {
		return nil, err
	}

	var doc bson.D
	bson.Unmarshal(data, &doc)

	return &doc, nil
}

func ParseToDatabaseId(ids ...string) ([]primitive.ObjectID, error) {
	objectIds := []primitive.ObjectID{}

	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, objectId)
	}

	return objectIds, nil
}
