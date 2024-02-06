package app_test

import (
	"errors"
	"log"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestGetUserById_Do(t *testing.T) {
	type MockStructure struct {
		Foo string
	}

	t.Run("should return error when failed to call database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		getUserByIdImpl := app.NewGetUserByIdImpl(mockCrudRepository, mockUserRepository)

		mockExpectedError := errors.New("something goes wrong")
		id := primitive.NewObjectID().Hex()

		mockCrudRepository.
			EXPECT().
			GetById(database.UsersCollection, id, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := getUserByIdImpl.Do(id)
		if err == nil {
			log.Fatal(err)
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		getUserByIdImpl := app.NewGetUserByIdImpl(mockCrudRepository, mockUserRepository)

		id := primitive.NewObjectID().Hex()

		mockCrudRepository.
			EXPECT().
			GetById(database.UsersCollection, id, gomock.Any()).
			Times(1).
			Return(nil)

		_, err := getUserByIdImpl.Do(id)
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
	})
}
