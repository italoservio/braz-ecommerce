package app

import (
	"errors"
	"os"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

type CreateUserInterface interface {
	Do(createUser *CreateUserInput) (*CreateUserOutput, error)
}

type CreateUserImpl struct {
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewCreateUserImpl(
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *CreateUserImpl {
	return &CreateUserImpl{crudRepository: cr, userRepository: ur}
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

func (gu *CreateUserImpl) Do(createUser *CreateUserInput) (*CreateUserOutput, error) {
	secret := os.Getenv("ENC_SECRET")
	encryptionData, err := encryption.Encrypt(secret, createUser.Password)

	if err != nil {
		return nil, errors.New(exception.CodeInternal)
	}

	id, err := gu.crudRepository.CreateOne(database.UsersCollection, &CreateUserDatabase{
		User: domain.User{
			Type:      createUser.Type,
			FirstName: createUser.FirstName,
			LastName:  createUser.LastName,
			Email:     createUser.Email,
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
