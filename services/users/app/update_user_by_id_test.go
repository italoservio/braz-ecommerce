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
	updateUserByIdImpl *app.UpdateUserByIdImpl
}

func BeforeEach_TestUpdateUserById(t *testing.T) *TestingDependencies_TestUpdateUser {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	encryption := mocks.NewMockEncryptionInterface(ctrl)
	mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
	mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

	updateUserByIdImpl := app.NewUpdateUserByIdImpl(
		encryption,
		mockCrudRepository,
		mockUserRepository,
	)

	return &TestingDependencies_TestUpdateUser{
		ctx:                ctx,
		ctrl:               ctrl,
		encryption:         encryption,
		mockCrudRepository: mockCrudRepository,
		mockUserRepository: mockUserRepository,
		updateUserByIdImpl: updateUserByIdImpl,
	}
}

func TestUpdateUser_Do(t *testing.T) {
	t.Run("should return error when failed to call database in GetByEmail", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUserById(t)
		defer deps.ctrl.Finish()

		mockEmail := "goo@gle.com"

		mockExpectedError := errors.New("something goes wrong")
		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, mockEmail, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.updateUserByIdImpl.Do(deps.ctx, id, &app.UpdateUserByIdInput{Email: mockEmail})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return a permission error because the id sent is different from the one found in the email", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUserById(t)
		defer deps.ctrl.Finish()

		mockEmail := "goo@gle.com"

		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, mockEmail, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				email string,
				structure *domain.UserDatabaseNoPassword,
			) error {
				*structure = domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: primitive.NewObjectID().Hex(),
					},
				}

				return nil
			})

		_, err := deps.updateUserByIdImpl.Do(deps.ctx, id, &app.UpdateUserByIdInput{Email: mockEmail})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
		assert.Equal(t, "EPERMISSION", err.Error(), "should return the expected error code")
	})

	t.Run("should return error when failed to call database in UpdateById", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUserById(t)
		defer deps.ctrl.Finish()

		mockEmail := "goo@gle.com"

		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, mockEmail, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				email string,
				structure *domain.UserDatabaseNoPassword,
			) error {
				*structure = domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: id,
					},
				}

				return nil
			})

		mockExpectedError := errors.New("something goes wrong")

		deps.mockCrudRepository.
			EXPECT().
			UpdateById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.updateUserByIdImpl.Do(deps.ctx, id, &app.UpdateUserByIdInput{Email: mockEmail})

		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return error when failed to call Encrypt", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUserById(t)
		defer deps.ctrl.Finish()

		mockEmail := "goo@gle.com"

		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, mockEmail, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				email string,
				structure *domain.UserDatabaseNoPassword,
			) error {
				*structure = domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: id,
					},
				}

				return nil
			})

		mockPassword := "test"
		mockExpectedError := errors.New("secret or text is empty")

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(nil, mockExpectedError)

		_, err := deps.updateUserByIdImpl.Do(deps.ctx, id, &app.UpdateUserByIdInput{Email: mockEmail, Password: "test"})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "something goes wrong")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		deps := BeforeEach_TestUpdateUserById(t)
		defer deps.ctrl.Finish()

		mockEmail := "goo@gle.com"

		id := primitive.NewObjectID().Hex()

		deps.mockUserRepository.
			EXPECT().
			GetByEmail(gomock.Any(), database.UsersCollection, mockEmail, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				email string,
				structure *domain.UserDatabaseNoPassword,
			) error {
				*structure = domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{
						Id: id,
					},
				}

				return nil
			})

		deps.mockCrudRepository.
			EXPECT().
			UpdateById(gomock.Any(), database.UsersCollection, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")
		mockPassword := "test"

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(&encryption.EncryptedText{EncryptedText: "", Salt: ""}, nil)

		_, err := deps.updateUserByIdImpl.Do(deps.ctx, id, &app.UpdateUserByIdInput{Email: mockEmail, Password: "test"})
		if err != nil {
			t.Fail()
		}

		assert.Nil(t, err, "should not return an error")
	})

}
