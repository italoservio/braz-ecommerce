package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	database_mocks "github.com/italoservio/braz_ecommerce/packages/database/mocks"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestDeleteUserById struct {
	ctx                context.Context
	ctrl               *gomock.Controller
	mockCrudRepository *database_mocks.MockCrudRepositoryInterface
	mockUserRepository *mocks.MockUserRepositoryInterface
	deleteUserByIdImpl *app.DeleteUserByIdImpl
}

func BeforeEach_TestDeleteUserById(t *testing.T) *TestingDependencies_TestDeleteUserById {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockCrudRepository := database_mocks.NewMockCrudRepositoryInterface(ctrl)
	mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

	deleteUserByIdImpl := app.NewDeleteUserByIdImpl(mockCrudRepository, mockUserRepository)

	return &TestingDependencies_TestDeleteUserById{
		ctx:                ctx,
		ctrl:               ctrl,
		mockCrudRepository: mockCrudRepository,
		mockUserRepository: mockUserRepository,
		deleteUserByIdImpl: deleteUserByIdImpl,
	}
}

func TestDeleteUserById_Do(t *testing.T) {
	t.Run("should return nil when deleted with success", func(t *testing.T) {
		deps := BeforeEach_TestDeleteUserById(t)

		id := primitive.NewObjectID().Hex()

		deps.mockCrudRepository.EXPECT().
			DeleteById(gomock.Any(), database.UsersCollection, id).
			Times(1).
			Return(nil)

		err := deps.deleteUserByIdImpl.Do(deps.ctx, id)

		assert.Nil(t, err, "should return nil")
	})

	t.Run("should return the error when failed to delete", func(t *testing.T) {
		deps := BeforeEach_TestDeleteUserById(t)

		id := primitive.NewObjectID().Hex()

		deps.mockCrudRepository.EXPECT().
			DeleteById(gomock.Any(), database.UsersCollection, id).
			Times(1).
			Return(errors.New(exception.CodeDatabaseFailed))

		err := deps.deleteUserByIdImpl.Do(deps.ctx, id)

		assert.NotNil(t, err, "should return the error")
	})
}
