package start

import (
	"fmt"
	"os"
)

type Env struct {
	PORT    string
	DB_URI  string
	DB_NAME string
}

func NewEnv() *Env {
	return &Env{
		PORT:    fmt.Sprintf(":%v", os.Getenv("PORT")),
		DB_URI:  os.Getenv("DB_URI"),
		DB_NAME: os.Getenv("DB_NAME"),
	}
}
