package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN string // e.g. "host=localhost user=postgres password=root dbname=freeconnect sslmode=disable"
	Port   string // e.g. "8080"
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists.
	_ = godotenv.Load()
	cfg := &Config{
		DB_DSN: "host=localhost user=postgres password=root dbname=freeconnect sslmode=disable",
		Port:   os.Getenv("PORT"),
	}
	if cfg.DB_DSN == "" {
		return nil, errors.New("DB_DSN is not set")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg, nil
}
