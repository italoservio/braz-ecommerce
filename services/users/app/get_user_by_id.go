package app

import (
	"fmt"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetUserByIdInterface interface {
	Do(id string) (*NewGetUserByIdOutput, error)
}

type GetUserByIdImpl struct {
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewGetUserByIdImpl(
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *GetUserByIdImpl {
	return &GetUserByIdImpl{crudRepository: cr, userRepository: ur}
}

type NewGetUserByIdOutput struct {
	*domain.User        `bson:",inline"`
	*domain.UserControl `bson:",inline"`
}

func (gu *GetUserByIdImpl) Do(id string) (*NewGetUserByIdOutput, error) {

	var output NewGetUserByIdOutput

	res, err := gu.crudRepository.GetById(database.UserCollection, id)
	if err != nil {
		return nil, err
	}

	if err := res.Decode(&output); err != nil {
		if err == mongo.ErrNoDocuments {
			// todo: should return default error
			return nil, fmt.Errorf("user with id '%s' not found", id)
		}
	}

	return &output, nil
}
