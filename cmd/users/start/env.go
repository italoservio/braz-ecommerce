package start

import (
	"fmt"
	"os"
)

type EnvironmentVariables struct {
	PORT       string
	DB_URI     string
	DB_NAME    string
	ENC_SECRET string
}

var Env *EnvironmentVariables

func NewEnv() *EnvironmentVariables {
	return &EnvironmentVariables{
		PORT:       fmt.Sprintf(":%v", os.Getenv("PORT")),
		DB_URI:     os.Getenv("DB_URI"),
		DB_NAME:    os.Getenv("DB_NAME"),
		ENC_SECRET: os.Getenv("ENC_SECRET"),
	}
}
