package app

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

type GetUserByIdInterface interface {
	Do(id string) string
}

type GetUserByIdImpl struct {
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewGetUserByIdImpl(
	ur storage.UserRepositoryInterface,
	cr database.CrudRepositoryInterface,
) *GetUserByIdImpl {
	return &GetUserByIdImpl{userRepository: ur}
}

func (gu *GetUserByIdImpl) Do(id string) string {
	return "sample"
}
