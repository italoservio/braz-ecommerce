package app

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
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

	err := gu.crudRepository.GetById(database.UsersCollection, id, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
