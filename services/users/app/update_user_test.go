package app_test

import (
	"errors"
	"os"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser_Do(t *testing.T) {

	t.Run("should return error when failed to call database UpdateById", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		updateUserImpl := app.NewUpdateUserImpl(mockCrudRepository, mockUserRepository)

		mockExpectedError := errors.New("something goes wrong")
		id := primitive.NewObjectID().Hex()

		mockCrudRepository.
			EXPECT().
			GetByEmail(database.UsersCollection, "teste", gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		mockCrudRepository.
			EXPECT().
			UpdateById(database.UsersCollection, id, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := updateUserImpl.Do(&app.UpdateUserInput{Email: "teste"}, id, app.UpdateUserOutput{})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return error when failed to call database GetById", func(t *testing.T) {
		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExpectedError := errors.New("something goes wrong")
		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		updateUserImpl := app.NewUpdateUserImpl(mockCrudRepository, mockUserRepository)

		mockCrudRepository.
			EXPECT().
			GetByEmail(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockCrudRepository.
			EXPECT().
			UpdateById(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockCrudRepository.
			EXPECT().
			GetById(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := updateUserImpl.Do(&app.UpdateUserInput{Email: "testeteste", Password: "123"}, "",
			app.UpdateUserOutput{
				UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: ""},
				},
			},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return a permission error because the id sent is different from the one found in the email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		updateUserImpl := app.NewUpdateUserImpl(mockCrudRepository, mockUserRepository)

		id := primitive.NewObjectID().Hex()

		mockCrudRepository.
			EXPECT().
			GetByEmail(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		_, err := updateUserImpl.Do(&app.UpdateUserInput{Email: "teste"}, id,
			app.UpdateUserOutput{
				UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: ""},
				},
			},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "EPERMISSION")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
		mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

		updateUserImpl := app.NewUpdateUserImpl(mockCrudRepository, mockUserRepository)

		mockCrudRepository.
			EXPECT().
			GetByEmail(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockCrudRepository.
			EXPECT().
			UpdateById(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockCrudRepository.
			EXPECT().
			GetById(database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		_, err := updateUserImpl.Do(&app.UpdateUserInput{Email: "testeteste", Password: "123"}, "",
			app.UpdateUserOutput{
				UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: ""},
				},
			},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "nil")
	})
}
