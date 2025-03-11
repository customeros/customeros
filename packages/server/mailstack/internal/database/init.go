package database

import (
	"log"

	"gorm.io/gorm"
)

func InitMailstackDatabase(dbConfig *DatabaseConfig) (*gorm.DB, error) {
	db, err := NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db, nil
}
