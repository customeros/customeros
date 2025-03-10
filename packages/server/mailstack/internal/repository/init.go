package repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/storage"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/config"
	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type Repositories struct {
	EmailRepository           interfaces.EmailRepository
	EmailAttachmentRepository interfaces.EmailAttachmentRepository
	MailboxRepository         interfaces.MailboxRepository
	MailboxSyncRepository     interfaces.MailboxSyncRepository
	MessageStateRepository    interfaces.MessageStateRepository
}

func InitRepositories(mailstackDB *gorm.DB, r2Config *config.R2StorageConfig) *Repositories {
	emailAttachmentStorage := storage.NewR2StorageService(
		r2Config.AccountID,
		r2Config.AccessKeyID,
		r2Config.AccessKeySecret,
		r2Config.EmailAttachmentBucket,
		false, // private access
	)

	return &Repositories{
		EmailRepository:           NewEmailRepository(mailstackDB),
		EmailAttachmentRepository: NewEmailAttachmentRepository(mailstackDB, emailAttachmentStorage),
		MailboxRepository:         NewMailboxRepository(mailstackDB),
		MailboxSyncRepository:     NewMailboxSyncRepository(mailstackDB),
		MessageStateRepository:    NewMessageStateRepository(mailstackDB),
	}
}

func MigrateDB(mailstackDB *gorm.DB) error {
	return mailstackDB.AutoMigrate(
		&models.Email{},
		&models.EmailAttachment{},
		&models.Mailbox{},
		&models.MailboxSyncState{},
		&models.MessageState{},
	)
}
