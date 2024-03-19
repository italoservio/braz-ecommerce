package app

import (
	"context"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

//go:generate mockgen --source=get_user_by_id.go --destination=../mocks/get_user_by_id_interface_mock.go --package=mocks
type GetUserByIdInterface interface {
	Do(ctx context.Context, input *GetUserByIdInput) (*GetUserByIdOutput, error)
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

type GetUserByIdInput struct {
	Id      string
	Deleted bool
}

type GetUserByIdOutput struct {
	*domain.UserDatabaseNoPassword `bson:",inline"`
}

func (gu *GetUserByIdImpl) Do(ctx context.Context, input *GetUserByIdInput) (*GetUserByIdOutput, error) {
	var output GetUserByIdOutput

	err := gu.crudRepository.GetById(ctx, database.UsersCollection, input.Id, input.Deleted, &output)

	if err != nil {
		return nil, err
	}

	return &output, nil
}
