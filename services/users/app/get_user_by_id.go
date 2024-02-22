package app

import (
	"context"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

type GetUserByIdInterface interface {
	Do(ctx context.Context, id string) (*GetUserByIdOutput, error)
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

type GetUserByIdOutput struct {
	*domain.UserDatabaseNoPassword `bson:",inline"`
}

func (gu *GetUserByIdImpl) Do(ctx context.Context, id string) (*GetUserByIdOutput, error) {
	var output GetUserByIdOutput

	err := gu.crudRepository.GetById(ctx, database.UsersCollection, id, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
