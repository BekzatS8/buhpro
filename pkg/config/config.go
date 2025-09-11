package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppAddr        string // e.g. ":8080"
	JWTSecret      string
	JTTTLMin       int // access token TTL in minutes
	RefreshTTLDays int // refresh token TTL in days
	DatabaseURL    string
	// add other fields you already have...
}

// Load reads config from environment variables with sensible defaults.
// You can adapt it to use viper/envconfig/whatever your project uses.
func Load() *Config {
	cfg := &Config{
		AppAddr:        getEnv("APP_ADDR", ":8080"),
		JWTSecret:      getEnv("JWT_SECRET", "replace-me-with-secure-secret"),
		JTTTLMin:       getEnvInt("JWT_TTL_MIN", 60),      // default 60 minutes
		RefreshTTLDays: getEnvInt("REFRESH_TTL_DAYS", 30), // default 30 days
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:1234@localhost:5432/buhpro?sslmode=disable"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
