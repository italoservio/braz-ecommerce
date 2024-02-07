package app_test

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser_Do(t *testing.T) {
	t.Run("should return error when failed to call database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		createUserImpl := app.NewCreateUserImpl(mockCrudRepository, mockUserRepository)

		mockExpectedError := errors.New("something goes wrong")

		mockCrudRepository.
			EXPECT().
			CreateOne(database.UsersCollection, gomock.Any()).
			Times(1).
			Return("", mockExpectedError)

		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		_, err := createUserImpl.Do(&app.CreateUserInput{Password: "test"})
		if err == nil {
			log.Fatal(err)
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return error when failed to encrypt password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		createUserImpl := app.NewCreateUserImpl(mockCrudRepository, mockUserRepository)

		_, err := createUserImpl.Do(&app.CreateUserInput{Password: ""})
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

		createUserImpl := app.NewCreateUserImpl(mockCrudRepository, mockUserRepository)

		mockCrudRepository.
			EXPECT().
			CreateOne(database.UsersCollection, gomock.Any()).
			Times(1).
			Return("", nil)

		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		_, err := createUserImpl.Do(&app.CreateUserInput{Password: "test"})
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
	})
}
