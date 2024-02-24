package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

type UpdateUserInterface interface {
	Do(ctx context.Context, updateUser *UpdateUserInput, id string, output UpdateUserOutput) (*UpdateUserOutput, error)
}

type UpdateUserImpl struct {
	encryption     encryption.EncryptionInterface
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewUpdateUserImpl(
	en encryption.EncryptionInterface,
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *UpdateUserImpl {
	return &UpdateUserImpl{
		encryption:     en,
		crudRepository: cr,
		userRepository: ur,
	}
}

type UpdateUserInput struct {
	FirstName string    `json:"first_name" validate:"omitempty,min=1,max=100" bson:"first_name,omitempty"`
	LastName  string    `json:"last_name" validate:"omitempty,min=1,max=100" bson:"last_name,omitempty"`
	Email     string    `json:"email" validate:"omitempty,min=1,max=100" bson:"email,omitempty"`
	Type      string    `json:"type" validate:"omitempty,min=1,max=100" bson:"type,omitempty"`
	Password  string    `json:"password" validate:"omitempty,min=1,max=100" bson:"password,omitempty"`
	UpdatedAt time.Time `json:"updated_at" validate:"omitempty" bson:"updated_at,omitempty"`
}

type UpdateUserOutput struct {
	*domain.UserDatabaseNoPassword `bson:",inline"`
}

func (gu *UpdateUserImpl) Do(ctx context.Context, updateUser *UpdateUserInput, id string, output UpdateUserOutput) (*UpdateUserOutput, error) {

	err := gu.crudRepository.GetByEmail(ctx, database.UsersCollection, updateUser.Email, &output)

	if err == nil && output.Id != id {
		return nil, errors.New(exception.CodePermission)
	}

	if updateUser.Password != "" {
		secret := os.Getenv("ENC_SECRET")
		encryptionData, err := gu.encryption.Encrypt(ctx, secret, updateUser.Password)

		fmt.Printf("%v", err)
		updateUser.Password = encryptionData.EncryptedText
		updateUser.Password = encryptionData.Salt
	}

	err = gu.crudRepository.UpdateById(ctx, database.UsersCollection, id, &updateUser)

	if err != nil {
		return nil, err
	}

	var userOutput UpdateUserOutput

	err = gu.crudRepository.GetById(ctx, database.UsersCollection, id, &userOutput)

	if err != nil {
		return nil, err
	}

	return &userOutput, nil
}
