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
	//Test comment
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Option A: Provide a fallback (localhost):
		// dsn = "host=localhost user=postgres password=root dbname=freeconnect sslmode=disable"

		// Option B: Or simply fail if it's not set:
		return nil, errors.New("DB_DSN is not set (and no fallback provided)")
	}

	// Read PORT from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg := &Config{
		DB_DSN: dsn,
		Port:   port,
	}
	return cfg, nil
}
