package repository

import (
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"gorm.io/gorm"
	"log"
)

type PostgresRepositories struct {
	ConversationEventRepository ConversationEventRepository
	CommonRepositories          *commonRepository.PostgresCommonRepositoryContainer
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		ConversationEventRepository: NewConversationEventRepository(db),
		CommonRepositories:          commonRepository.InitCommonRepositories(db),
	}

	err := db.AutoMigrate(&entity.ConversationEvent{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return p
}
