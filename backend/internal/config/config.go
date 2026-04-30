package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiryHour time.Duration
}

func Load() (Config, error) {
	expiryRaw := getEnv("JWT_EXPIRY_HOURS", "24")
	expiryInt, err := strconv.Atoi(expiryRaw)
	if err != nil || expiryInt <= 0 {
		return Config{}, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %s", expiryRaw)
	}

	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/product_management?sslmode=disable")
	secret := getEnv("JWT_SECRET", "super-secret-dev-key")

	if len(secret) < 10 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 10 characters")
	}

	return Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   dbURL,
		JWTSecret:     secret,
		JWTExpiryHour: time.Duration(expiryInt) * time.Hour,
	}, nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
