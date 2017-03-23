package server

import (
	"os"

	"github.com/joho/godotenv"
)

func (a API) SetupViper() {
	var filename string
	env := os.Getenv("BASEAPI_ENV")

	if env == "testing" {
		filename = "../.env.testing"
	} else if env == "prod" {
		filename = ".env.prod"
	} else {
		filename = ".env"
	}

	err := godotenv.Overload(filename)
	if err != nil {
		panic(err)
	}

	a.Config.SetEnvPrefix("baseapi")
	a.Config.AutomaticEnv()
}
