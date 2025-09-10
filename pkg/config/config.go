package config

import (
	"os"
)

type Config struct {
	AppEnv  string
	AppAddr string
	DB_DSN  string
}

func Load() *Config {
	return &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppAddr: getEnv("APP_ADDR", ":8080"),
		DB_DSN:  getEnv("DB_DSN", "postgres://postgres:1234@localhost:5432/buhpro?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
