package app_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestCreateUser struct {
	ctx                context.Context
	ctrl               *gomock.Controller
	encryption         *mocks.MockEncryptionInterface
	mockCrudRepository *mocks.MockCrudRepositoryInterface
	mockUserRepository *mocks.MockUserRepositoryInterface
	createUserImpl     *app.CreateUserImpl
}

func BeforeEach_TestCreateUser(t *testing.T) *TestingDependencies_TestCreateUser {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	encryption := mocks.NewMockEncryptionInterface(ctrl)
	mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)
	mockUserRepository := mocks.NewMockUserRepositoryInterface(ctrl)

	createUserImpl := app.NewCreateUserImpl(encryption, mockCrudRepository, mockUserRepository)

	return &TestingDependencies_TestCreateUser{
		ctx:                ctx,
		ctrl:               ctrl,
		encryption:         encryption,
		mockCrudRepository: mockCrudRepository,
		mockUserRepository: mockUserRepository,
		createUserImpl:     createUserImpl,
	}
}

func TestCreateUser_Do(t *testing.T) {
	t.Run("should return error when failed to call database", func(t *testing.T) {
		deps := BeforeEach_TestCreateUser(t)
		defer deps.ctrl.Finish()

		mockExpectedError := errors.New("something goes wrong")
		mockPassword := "test"

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(&encryption.EncryptedText{EncryptedText: "", Salt: ""}, nil)

		deps.mockCrudRepository.
			EXPECT().
			CreateOne(gomock.Any(), database.UsersCollection, gomock.Any()).
			Times(1).
			Return("", mockExpectedError)

		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		_, err := deps.createUserImpl.Do(deps.ctx, &app.CreateUserInput{Password: mockPassword})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return error when failed to encrypt password", func(t *testing.T) {
		deps := BeforeEach_TestCreateUser(t)
		defer deps.ctrl.Finish()

		mockExpectedError := errors.New("secret or text is empty")
		mockPassword := ""

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(nil, mockExpectedError)

		_, err := deps.createUserImpl.Do(deps.ctx, &app.CreateUserInput{Password: mockPassword})
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		deps := BeforeEach_TestCreateUser(t)
		defer deps.ctrl.Finish()

		mockPassword := "test"

		deps.encryption.
			EXPECT().
			Encrypt(gomock.Any(), gomock.Any(), mockPassword).
			Times(1).
			Return(&encryption.EncryptedText{EncryptedText: "", Salt: ""}, nil)

		deps.mockCrudRepository.
			EXPECT().
			CreateOne(gomock.Any(), database.UsersCollection, gomock.Any()).
			Times(1).
			Return("", nil)

		os.Setenv("ENC_SECRET", "2zmXvZa93wneR1w1L63i9cAUzSIzPdd6")

		_, err := deps.createUserImpl.Do(deps.ctx, &app.CreateUserInput{Password: mockPassword})
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
	})
}
