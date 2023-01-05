package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositories struct {
	ConversationEventRepository ConversationEventRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		ConversationEventRepository: NewConversationEventRepository(db),
	}

	err := db.AutoMigrate(&entity.ConversationEvent{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
