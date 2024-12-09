package config

import (
	"fmt"
	"os"
)

type Config struct {
	DSN string
}

func LoadConfig() *Config {
	return &Config{
		DSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	}
}
