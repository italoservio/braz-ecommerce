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

type UpdateUserByIdInterface interface {
	Do(ctx context.Context, id string, input *UpdateUserByIdInput) (*UpdateUserByIdOutput, error)
}

type UpdateUserByIdImpl struct {
	encryption     encryption.EncryptionInterface
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewUpdateUserByIdImpl(
	en encryption.EncryptionInterface,
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *UpdateUserByIdImpl {
	return &UpdateUserByIdImpl{
		encryption:     en,
		crudRepository: cr,
		userRepository: ur,
	}
}

type UpdateUserByIdInput struct {
	FirstName string    `json:"first_name" validate:"omitempty,min=1,max=100" bson:"first_name,omitempty"`
	LastName  string    `json:"last_name" validate:"omitempty,min=1,max=100" bson:"last_name,omitempty"`
	Email     string    `json:"email" validate:"omitempty,min=1,max=100" bson:"email,omitempty"`
	Type      string    `json:"type" validate:"omitempty,min=1,max=100" bson:"type,omitempty"`
	Password  string    `json:"password" validate:"omitempty,min=1,max=100" bson:"password,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

type UpdateUserByIdOutput struct {
	*domain.UserDatabaseNoPassword `bson:",inline"`
}

func (gu *UpdateUserByIdImpl) Do(
	ctx context.Context,
	id string,
	input *UpdateUserByIdInput,
) (*UpdateUserByIdOutput, error) {
	var existentUser domain.UserDatabaseNoPassword

	if input.Email != "" {
		err := gu.userRepository.GetByEmail(
			ctx,
			database.UsersCollection,
			input.Email,
			&existentUser,
		)

		if err != nil {
			return nil, err
		}

		if existentUser != (domain.UserDatabaseNoPassword{}) && existentUser.Id != id {
			return nil, errors.New(exception.CodePermission)
		}
	}

	if input.Password != "" {
		secret := os.Getenv("ENC_SECRET")
		encryptionData, err := gu.encryption.Encrypt(ctx, secret, input.Password)

		if err != nil {
			return nil, err
		}

		input.Password = encryptionData.EncryptedText
		input.Password = encryptionData.Salt
	}

	input.UpdatedAt = time.Now()

	var output = UpdateUserByIdOutput{}

	err := gu.crudRepository.UpdateById(ctx, database.UsersCollection, id, &input, &output)

	if err != nil {
		return nil, err
	}

	return &output, nil
}
