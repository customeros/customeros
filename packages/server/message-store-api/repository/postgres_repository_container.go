package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositories struct {
	ConversationEventRepository ConversationEventRepository
	CommonRepositories          *commonRepository.Repositories
}

func InitRepositories(db *gorm.DB, driver *neo4j.Driver) *PostgresRepositories {
	p := &PostgresRepositories{
		ConversationEventRepository: NewConversationEventRepository(db),
		CommonRepositories:          commonRepository.InitRepositories(db, driver),
	}

	err := db.AutoMigrate(&entity.ConversationEvent{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
