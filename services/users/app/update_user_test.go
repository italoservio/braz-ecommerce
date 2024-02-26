package app_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestUpdateUser struct {
	ctx                context.Context
	ctrl               *gomock.Controller
	encryption         *mocks.MockEncryptionInterface
	mockCrudRepository *mocks.MockCrudRepositoryInterface
	mockUserRepository *mocks.MockUserRepositoryInterface
	updateUserImpl     *app.UpdateUserByIdImpl
}

func BeforeEach_TestUpdateUser(t *testing.T) *TestingDependencies_TestUpdateUser {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	encryption := mocks.NewMockEncryptionInterface(ctrl)
	mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
	mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

	updateUserImpl := app.NewUpdateUserImpl(encryption, mockCrudRepository, mockUserRepository)

	return &TestingDependencies_TestUpdateUser{
		ctx:                ctx,
		ctrl:               ctrl,
		encryption:         encryption,
		mockCrudRepository: mockCrudRepository,
		mockUserRepository: mockUserRepository,
		updateUserImpl:     updateUserImpl,
	}
}

func TestUpdateUser_Do(t *testing.T) {

	t.Run("should return error when failed to call database UpdateById", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUser(t)
		defer deps.ctrl.Finish()

		mockExpectedError := errors.New("something goes wrong")
		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, "teste", gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		deps.mockCrudRepository.
			EXPECT().
			UpdateById(gomock.Any(), database.UsersCollection, id, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.updateUserImpl.Do(deps.ctx, &app.UpdateUserInput{Email: "teste"}, id, app.UpdateUserOutput{})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return error when failed to call database GetById", func(t *testing.T) {
		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		deps := BeforeEach_TestUpdateUser(t)
		defer deps.ctrl.Finish()

		mockExpectedError := errors.New("something goes wrong")
		mockPassword := "test"

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(&encryption.EncryptedText{EncryptedText: "", Salt: ""}, nil)

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		deps.mockCrudRepository.
			EXPECT().
			UpdateById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		deps.mockCrudRepository.
			EXPECT().
			GetById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.updateUserImpl.Do(deps.ctx, &app.UpdateUserInput{Email: "testeteste", Password: mockPassword}, "",
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
		deps := BeforeEach_TestUpdateUser(t)
		defer deps.ctrl.Finish()

		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		_, err := deps.updateUserImpl.Do(deps.ctx, &app.UpdateUserInput{Email: "teste"}, id,
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

		deps := BeforeEach_TestUpdateUser(t)
		mockPassword := "test"

		defer deps.ctrl.Finish()

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(&encryption.EncryptedText{EncryptedText: "", Salt: ""}, nil)
		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		deps.mockCrudRepository.
			EXPECT().
			UpdateById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		deps.mockCrudRepository.
			EXPECT().
			GetById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		_, err := deps.updateUserImpl.Do(deps.ctx, &app.UpdateUserInput{Email: "testeteste", Password: mockPassword}, "",
			app.UpdateUserOutput{
				UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: ""},
				},
			},
		)
		if err != nil {
			t.Fail()
		}

		assert.Nil(t, err, "should not return an error")
	})
}
