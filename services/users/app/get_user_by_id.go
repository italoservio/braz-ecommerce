package app

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
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

	res, err := gu.crudRepository.GetById(database.UsersCollection, id)
	if err != nil {
		return nil, err
	}

	if err := res.Decode(&output); err != nil {
		if err == mongo.ErrNoDocuments {
			slog.Error(fmt.Sprintf("Unable to find user with id '%s'", id))

			return nil, errors.New(exception.CodeNotFound)
		}
	}

	return &output, nil
}
