//go:build dev

package config

import "github.com/joho/godotenv"

func init() {
	_ = godotenv.Load("../.env")
}
