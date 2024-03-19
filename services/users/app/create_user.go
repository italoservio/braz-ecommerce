package app

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

//go:generate mockgen --source=create_user.go --destination=../mocks/create_user_interface_mock.go --package=mocks
type CreateUserInterface interface {
	Do(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error)
}

type CreateUserImpl struct {
	encryption     encryption.EncryptionInterface
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewCreateUserImpl(
	en encryption.EncryptionInterface,
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *CreateUserImpl {
	return &CreateUserImpl{
		encryption:     en,
		crudRepository: cr,
		userRepository: ur,
	}
}

type CreateUserOutput struct {
	*database.DatabaseIdentifier
}

type CreateUserInput struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Email     string `json:"email" validate:"required,min=1,max=100"`
	Type      string `json:"type" validate:"required,min=5,max=100"`
	Password  string `json:"password" validate:"required,min=5,max=100"`
}

type CreateUserDatabase struct {
	domain.User                `bson:",inline"`
	domain.UserPassword        `bson:",inline"`
	database.DatabaseTimestamp `bson:",inline"`
}

func (gu *CreateUserImpl) Do(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	secret := os.Getenv("ENC_SECRET")
	encryptionData, err := gu.encryption.Encrypt(ctx, secret, input.Password)

	if err != nil {
		return nil, errors.New(exception.CodeInternal)
	}

	var existentUser domain.UserDatabaseNoPassword

	err = gu.userRepository.GetByEmail(
		ctx,
		database.UsersCollection,
		input.Email,
		&existentUser,
	)

	if err != nil {
		return nil, err
	}

	if existentUser != (domain.UserDatabaseNoPassword{}) {
		return nil, errors.New(exception.CodePermission)
	}

	id, err := gu.crudRepository.CreateOne(ctx, database.UsersCollection, &CreateUserDatabase{
		User: domain.User{
			Type:      input.Type,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Addresses: []domain.UserAddress{},
		},
		UserPassword: domain.UserPassword{
			Password:  encryptionData.EncryptedText,
			CipherKey: encryptionData.Salt,
		},
		DatabaseTimestamp: database.DatabaseTimestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	})

	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{DatabaseIdentifier: &database.DatabaseIdentifier{Id: id}}, nil
}
