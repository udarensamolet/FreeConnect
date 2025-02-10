package tests

import (
	"log"
	"os"

	"FreeConnect/internal/config"
	"FreeConnect/internal/models"
	"gorm.io/gorm"
)

// SetupTestDB initializes a connection to the freeconnect_test DB
func SetupTestDB() *gorm.DB {
	// 1. Load config (which may read .env)
	cfg, err := config.LoadTestConfig()
	if err != nil {
		log.Fatalf("Failed to load config for tests: %v", err)
	}

	// 2. Override the DB_DSN with your TEST_DB_DSN if present
	testDSN := os.Getenv("TEST_DB_DSN")
	if testDSN != "" {
		cfg.DB_DSN = testDSN
	}

	// 3. Connect
	db, err := models.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	// 4. AutoMigrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Skill{},
		&models.Proposal{},
		&models.Review{},
		&models.Transaction{},
		&models.Task{},
		&models.Notification{},
		&models.Invoice{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test DB: %v", err)
	}

	return db
}
