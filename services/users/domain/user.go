package domain

import "time"

type User struct {
	Type      string        `json:"type" bson:"type"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email" bson:"email"`
	Addresses []UserAddress `json:"addresses" bson:"addresses"`
}

type UserPassword struct {
	Password  string `json:"password" bson:"password"`
	cipherKey string `bson:"cipher_key"`
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

type UserControl struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

type UserDatabase struct {
	User
	UserPassword
	UserControl
}
