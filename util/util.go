package util

import (
	"os"

	"github.com/joho/godotenv"
)

func IsProd() bool {
	return os.Getenv("PRODUCTION") == "true"
}

func IsDev() bool {
	return !IsProd()
}

func AppEnv() string {
	if IsProd() {
		return "production"
	}
	return "development"
}

func LoadEnvVars() error {
	if err := godotenv.Load(".env.local"); err != nil {
		return err
	}
	return nil
}
