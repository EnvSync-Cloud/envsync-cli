package helper

import (
	"github.com/joho/godotenv"
)

func WriteEnv(env map[string]string) error {
	return godotenv.Write(env, ".env")
}
