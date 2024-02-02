package create_user

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct{ database *database.Database }

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{database: db}
}

type CrudRepositoryInterface interface {
	CreateUserRepository(payload *DTOCreateUserReq) (*mongo.SingleResult, error)
}

func (ur *UserRepository) CreateUser(payload *DTOCreateUserReq) error {
	return nil
}
