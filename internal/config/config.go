package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort               string
	DatabaseURL           string
	JWTSecret             string
	JWTIssuer             string
	JWTAccessTokenMinutes time.Duration
}

func Load() (Config, error) {
	accessTokenMinutes, err := getEnvAsInt("JWT_ACCESS_TOKEN_MINUTES", 15)
	if err != nil {
		return Config{}, fmt.Errorf("JWT_ACCESS_TOKEN_MINUTES must be a valid integer")
	}

	cfg := Config{
		AppPort:               getEnv("APP_PORT", "8080"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		JWTSecret:             getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTIssuer:             getEnv("JWT_ISSUER", "secret-shop"),
		JWTAccessTokenMinutes: time.Duration(accessTokenMinutes) * time.Minute,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvAsInt(key string, fallback int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}
