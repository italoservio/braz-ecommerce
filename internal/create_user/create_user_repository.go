package create_user

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
)

type UserRepository struct {
	crudRepository database.CrudRepositoryInterface
	database       *database.Database
}

func NewUserRepository(db *database.Database, crud database.CrudRepositoryInterface) *UserRepository {
	return &UserRepository{database: db, crudRepository: crud}
}

type CreateUserRepositoryInterface interface {
	CreateUser(payload *DTOCreateUserReq) error
}

func (ur *UserRepository) CreateUser(payload *DTOCreateUserReq) error {
	return nil
}
