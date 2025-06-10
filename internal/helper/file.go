package helper

import (
	"os"

	"github.com/joho/godotenv"
)

func WriteEnv(env map[string]string) error {
	return godotenv.Write(env, ".env")
}

func ReadEnv() (map[string]string, error) {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// Return empty map if .env file doesn't exist
		return make(map[string]string), nil
	}

	return godotenv.Read(".env")
}
