package database

import (
	"log"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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

func InitOpenlineDatabase(dbConfig *DatabaseConfig) (*gorm.DB, error) {
	db, err := NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(
		&postgres_entity.TenantWebhookApiKey{},
		&postgres_entity.AppKey{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db, nil
}
