package database

import (
	"log"
)

func InitDatabase(dbConfig *DatabaseConfig) (*DatabaseConnection, error) {
	db, err := NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db, nil
}
