package app

import (
	"errors"
	"os"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

type UpdateUserInterface interface {
	Do(updateUser *UpdateUserInput, user *NewGetUserByIdOutput) (*UpdateUserOutput, error)
}

type UpdateUserImpl struct {
	crudRepository database.CrudRepositoryInterface
	userRepository storage.UserRepositoryInterface
}

func NewUpdateUserImpl(
	cr database.CrudRepositoryInterface,
	ur storage.UserRepositoryInterface,
) *UpdateUserImpl {
	return &UpdateUserImpl{crudRepository: cr, userRepository: ur}
}

type UpdateUserOutput struct {
	*database.DatabaseIdentifier `bson:",inline"`
	*domain.User                 `bson:",inline"`
	*database.DatabaseTimestamp  `bson:",inline"`
}

type UpdateUserInput struct {
	Id        *string     `json:"id"`
	FirstName *string     `json:"first_name,omitempty" validate:"min=1,max=100" bson:"first_name,omitempty"`
	LastName  *string     `json:"last_name,omitempty" validate:"min=1,max=100" bson:"last_name,omitempty"`
	Email     *string     `json:"email,omitempty" validate:"min=1,max=100" bson:"email,omitempty"`
	Type      *string     `json:"type,omitempty" validate:"min=1,max=100" bson:"type,omitempty"`
	Password  *string     `json:"password,omitempty" validate:"min=1,max=100" bson:"password,omitempty"`
	Addresses []Addresses `json:"tag_list"`
}

type Addresses struct {
	Cep          *string `json:"cep,omitempty" validate:"required,min=1,max=100,string" bson:"cep,omitempty"`
	Street       *string `json:"street,omitempty" validate:"min=1,max=100" bson:"street,omitempty"`
	Neighborhood *string `json:"neighborhood,omitempty" validate:"min=1,max=100" bson:"neighborhood,omitempty"`
	State        *string `json:"state,omitempty" validate:"min=1,max=100" bson:"state,omitempty"`
	Country      *string `json:"country,omitempty" validate:"min=1,max=100" bson:"country,omitempty"`
	Number       *string `json:"number,omitempty" validate:"min=1,max=100" bson:"number,omitempty"`
	Complement   *string `json:"complement,omitempty" validate:"min=1,max=100" bson:"complement,omitempty"`
}

func (gu *UpdateUserImpl) Do(updateUser *UpdateUserInput, user *NewGetUserByIdOutput) (*UpdateUserOutput, error) {

	if *updateUser.Password != "" {
		secret := os.Getenv("ENC_SECRET")
		encryptionData, err := encryption.Encrypt(secret, *updateUser.Password)

		if err != nil {
			return nil, errors.New(exception.CodeInternal)
		}

		*updateUser.Password = encryptionData.EncryptedText
		*updateUser.Password = encryptionData.Salt
	}

	err := gu.crudRepository.UpdateById(database.UsersCollection, *updateUser.Id, &updateUser)

	if err != nil {
		return nil, err
	}
	return &UpdateUserOutput{}, nil
}
