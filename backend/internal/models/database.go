package models

import (
	"FreeConnect/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase opens a PostgreSQL connection using the DSN from the config.
func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB_DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
