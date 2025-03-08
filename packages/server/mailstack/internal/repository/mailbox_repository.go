package repository

import (
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type mailboxRepository struct {
	db *gorm.DB
}

func NewMailboxRepository(db *gorm.DB) interfaces.MailboxRepository {
	return &mailboxRepository{db: db}
}

func (r *mailboxRepository) GetMailboxes() ([]interfaces.MailboxConfig, error) {
	var mailboxes []models.Mailbox
	err := r.db.Find(&mailboxes).Error
	if err != nil {
		return nil, err
	}

	result := make([]interfaces.MailboxConfig, len(mailboxes))
	for i, m := range mailboxes {
		result[i] = interfaces.MailboxConfig{
			ID:       m.ID,
			Server:   m.Server,
			Port:     m.Port,
			Username: m.Username,
			Password: m.Password,
			Folders:  m.Folders,
			TLS:      m.TLS,
		}
	}

	return result, nil
}

func (r *mailboxRepository) GetMailbox(id string) (interfaces.MailboxConfig, error) {
	var mailbox models.Mailbox
	err := r.db.First(&mailbox, "id = ?", id).Error
	if err != nil {
		return interfaces.MailboxConfig{}, err
	}

	return interfaces.MailboxConfig{
		ID:       mailbox.ID,
		Server:   mailbox.Server,
		Port:     mailbox.Port,
		Username: mailbox.Username,
		Password: mailbox.Password,
		Folders:  mailbox.Folders,
		TLS:      mailbox.TLS,
	}, nil
}

func (r *mailboxRepository) SaveMailbox(config interfaces.MailboxConfig) error {
	mailbox := models.Mailbox{
		ID:       config.ID,
		Server:   config.Server,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Folders:  config.Folders,
		TLS:      config.TLS,
	}

	return r.db.Save(&mailbox).Error
}

func (r *mailboxRepository) DeleteMailbox(id string) error {
	return r.db.Delete(&models.Mailbox{}, "id = ?", id).Error
}
