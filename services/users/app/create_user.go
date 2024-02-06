package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/database"
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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
	Password  string `json:"password"`
}

type CreateUserDatabase struct {
	domain.User                `bson:",inline"`
	domain.UserPassword        `bson:",inline"`
	database.DatabaseTimestamp `bson:",inline"`
}

func (gu *CreateUserImpl) Do(createUser *CreateUserInput) (*CreateUserOutput, error) {

	password, cipherKey := EncryptPassword(createUser.Password)

	id, err := gu.crudRepository.CreateOne(database.UsersCollection, &CreateUserDatabase{
		User: domain.User{
			Type:      createUser.Type,
			FirstName: createUser.FirstName,
			LastName:  createUser.LastName,
			Email:     createUser.Email,
			Addresses: []domain.UserAddress{},
		},
		UserPassword: domain.UserPassword{
			Password:  password,
			CipherKey: cipherKey,
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

func EncryptPassword(plaintext string) (string, string) {

	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	salt := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, salt, []byte(plaintext), nil)

	return hex.EncodeToString(ciphertext), hex.EncodeToString(salt)
}
