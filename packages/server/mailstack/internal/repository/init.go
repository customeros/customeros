package repository

import (
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

type Repositories struct {
	MailboxRepository      interfaces.MailboxRepository
	MessageStateRepository interfaces.MessageStateRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		MailboxRepository:      NewMailboxRepository(db),
		MessageStateRepository: NewMessageStateRepository(db),
	}
}
