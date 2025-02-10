package config

import (
	"errors"
	"fmt"
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

// LoadTestConfig is used only in tests to load TEST_DB_DSN from .env or environment
func LoadTestConfig() (*Config, error) {
	// Make sure this actually loads .env:
	err := godotenv.Load("../.env")
	if err != nil {
		// Not always fatal if .env is absent, but let's debug:
		fmt.Printf("DEBUG: could not load .env: %v\n", err)
	}

	testDSN := os.Getenv("TEST_DB_DSN")
	if testDSN == "" {
		return nil, errors.New("TEST_DB_DSN not set (and no fallback provided)")
	}

	return &Config{DB_DSN: testDSN}, nil
}
