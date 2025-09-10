package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv    string
	AppAddr   string
	DB_DSN    string
	JWTSecret string
	JTTTLMin  int
}

func Load() *Config {
	ttl := 60
	if v := os.Getenv("JWT_TTL_MINUTES"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			ttl = i
		}
	}
	return &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppAddr:   getEnv("APP_ADDR", ":8080"),
		DB_DSN:    getEnv("DB_DSN", "postgres://postgres:1234@localhost:5432/buhpro?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "change_this_secret"),
		JTTTLMin:  ttl,
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
