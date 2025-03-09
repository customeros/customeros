package repository

import (
	common_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

type Repositories struct {
	MailboxRepository             interfaces.MailboxRepository
	MessageStateRepository        interfaces.MessageStateRepository
	TenantWebhookAPIKeyRepository common_repository.TenantWebhookApiKeyRepository
	AppKeyRepository              common_repository.AppKeyRepository
}

func InitRepositories(mailstackDB, commonDB *gorm.DB) *Repositories {
	return &Repositories{
		MailboxRepository:             NewMailboxRepository(mailstackDB),
		MessageStateRepository:        NewMessageStateRepository(mailstackDB),
		TenantWebhookAPIKeyRepository: common_repository.NewTenantWebhookApiKeyRepository(commonDB),
		AppKeyRepository:              common_repository.NewAppKeyRepo(commonDB),
	}
}
