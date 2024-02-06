package database

import (
	"go.mongodb.org/mongo-driver/bson"
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
