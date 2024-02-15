package domain

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
)

type User struct {
	Type      string        `json:"type" bson:"type"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email" bson:"email"`
	Addresses []UserAddress `json:"addresses" bson:"addresses"`
}

type UserPassword struct {
	Password  string `json:"password" bson:"password"`
	CipherKey string `bson:"cipher_key"`
}

type UserAddress struct {
	Cep          string  `json:"cep" bson:"cep"`
	Street       string  `json:"street" bson:"street"`
	Neighborhood string  `json:"neighborhood" bson:"neighborhood"`
	State        string  `json:"state" bson:"state"`
	Country      string  `json:"country" bson:"country"`
	Number       string  `json:"number" bson:"number"`
	Complement   *string `json:"complement" bson:"complement"`
}

type UserDatabase struct {
	User
	UserPassword
	database.DatabaseIdentifier
	database.DatabaseTimestamp
}
