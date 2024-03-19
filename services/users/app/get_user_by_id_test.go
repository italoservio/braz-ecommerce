package app_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	database_mocks "github.com/italoservio/braz_ecommerce/packages/database/mocks"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestGetUserById struct {
	ctx                context.Context
	ctrl               *gomock.Controller
	mockCrudRepository *database_mocks.MockCrudRepositoryInterface
	mockUserRepository *mocks.MockUserRepositoryInterface
	getUserByIdImpl    *app.GetUserByIdImpl
}

func BeforeEach_TestGetUserById(t *testing.T) *TestingDependencies_TestGetUserById {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockCrudRepository := database_mocks.NewMockCrudRepositoryInterface(ctrl)
	mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

	getUserByIdImpl := app.NewGetUserByIdImpl(mockCrudRepository, mockUserRepository)

	return &TestingDependencies_TestGetUserById{
		ctx:                ctx,
		ctrl:               ctrl,
		mockCrudRepository: mockCrudRepository,
		mockUserRepository: mockUserRepository,
		getUserByIdImpl:    getUserByIdImpl,
	}
}

func TestGetUserById_Do(t *testing.T) {
	type MockStructure struct {
		Foo string
	}

	t.Run("should return error when failed to call database", func(t *testing.T) {
		deps := BeforeEach_TestGetUserById(t)

		mockExpectedError := errors.New("something goes wrong")
		id := primitive.NewObjectID().Hex()

		deps.mockCrudRepository.
			EXPECT().
			GetById(gomock.Any(), database.UsersCollection, id, false, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.getUserByIdImpl.Do(deps.ctx, &app.GetUserByIdInput{Id: id, Deleted: false})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		deps := BeforeEach_TestGetUserById(t)
		id := primitive.NewObjectID().Hex()

		deps.mockCrudRepository.
			EXPECT().
			GetById(gomock.Any(), database.UsersCollection, id, false, gomock.Any()).
			Times(1).
			Return(nil)

		_, err := deps.getUserByIdImpl.Do(deps.ctx, &app.GetUserByIdInput{Id: id, Deleted: false})
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
	})
}
