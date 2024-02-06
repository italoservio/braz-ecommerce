package storage

import "github.com/italoservio/braz_ecommerce/packages/database"

type UserRepositoryInterface interface{}

type UserRepositoryImpl struct {
	database *database.Database
}

func NewUserRepositoryImpl(db *database.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{database: db}
}
