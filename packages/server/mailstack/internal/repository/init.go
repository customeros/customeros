package repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/storage"
	common_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/config"
	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

type Repositories struct {
	AppKeyRepository              common_repository.AppKeyRepository
	EmailRepository               interfaces.EmailRepository
	EmailAttachmentRepository     interfaces.EmailAttachmentRepository
	MailboxRepository             interfaces.MailboxRepository
	MessageStateRepository        interfaces.MessageStateRepository
	TenantWebhookAPIKeyRepository common_repository.TenantWebhookApiKeyRepository
}

func InitRepositories(mailstackDB, commonDB *gorm.DB, r2Config *config.R2StorageConfig) *Repositories {
	emailAttachmentStorage := storage.NewR2StorageService(
		r2Config.AccountID,
		r2Config.AccessKeyID,
		r2Config.AccessKeySecret,
		r2Config.EmailAttachmentBucket,
		false, // private access
	)

	return &Repositories{
		AppKeyRepository:              common_repository.NewAppKeyRepo(commonDB),
		EmailRepository:               NewEmailRepository(mailstackDB),
		EmailAttachmentRepository:     NewEmailAttachmentRepository(mailstackDB, emailAttachmentStorage),
		MailboxRepository:             NewMailboxRepository(mailstackDB),
		MessageStateRepository:        NewMessageStateRepository(mailstackDB),
		TenantWebhookAPIKeyRepository: common_repository.NewTenantWebhookApiKeyRepository(commonDB),
	}
}
