package app

import (
	"context"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

//go:generate mockgen --source=delete_user_by_id.go --destination=../mocks/delete_user_by_id_interface_mock.go --package=mocks
type DeleteUserByIdInterface interface {
	Do(ctx context.Context, id string) error
}

type DeleteUserByIdImpl struct {
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewDeleteUserByIdImpl(
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *DeleteUserByIdImpl {
	return &DeleteUserByIdImpl{crudRepository: cr, userRepository: ur}
}

func (gu *DeleteUserByIdImpl) Do(ctx context.Context, id string) error {
	err := gu.crudRepository.DeleteById(ctx, database.UsersCollection, id)
	if err != nil {
		return err
	}

	return nil
}
