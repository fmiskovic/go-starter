package util

import (
	"os"
	"strings"

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

func GetEnvOrDefault(key string, def string) string {
	env, ok := os.LookupEnv(key)
	if ok && IsNotBlank(env) {
		return env
	}
	return def
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsNotEmpty(s string) bool {
	return len(s) > 0
}

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsNotBlank(s string) bool {
	return !IsBlank(s)
}
