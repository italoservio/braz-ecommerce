package database

import "time"

type DatabaseIdentifier struct {
	Id string `json:"id" bson:"_id"`
}

type DatabaseTimestamp struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}
