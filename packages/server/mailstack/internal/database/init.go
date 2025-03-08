package database

import (
	"log"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/config"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

func InitDatabase(dbConfig *config.MailstackDatabaseConfig) (*gorm.DB, error) {
	db, err := NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(
		&models.Mailbox{},
		&models.MessageState{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db, nil
}
