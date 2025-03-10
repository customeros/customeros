package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Mailbox struct {
	ID       string             `gorm:"column:id;type:varchar(50);primaryKey"`
	Provider enum.EmailProvider `gorm:"column:provider;type:varchar(50);index;not null"`

	// IMAP Configuration
	ImapServer   string `gorm:"column:imap_server;type:varchar(255);not null"`
	ImapPort     int    `gorm:"column:imap_port;not null"`
	ImapUsername string `gorm:"column:imap_username;type:varchar(255);not null"`
	ImapPassword string `gorm:"column:imap_password;type:varchar(255);not null"`
	ImapTLS      bool   `gorm:"column:imap_tls;not null;default:true"`

	// SMTP Configuration
	SmtpServer   string `gorm:"column:smtp_server;type:varchar(255);not null"`
	SmtpPort     int    `gorm:"column:smtp_port;not null"`
	SmtpUsername string `gorm:"column:smtp_username;type:varchar(255);not null"`
	SmtpPassword string `gorm:"column:smtp_password;type:varchar(255);not null"`
	SmtpTLS      bool   `gorm:"column:smtp_tls;not null;default:true"`

	// Other Configuration
	Folders      pq.StringArray `gorm:"column:folders;type:text[];not null"`
	DisplayName  string         `gorm:"column:display_name;type:varchar(255)"`
	EmailAddress string         `gorm:"column:email_address;type:varchar(255);index"`

	// Status Information
	LastSynced   *time.Time `gorm:"column:last_synced;type:timestamp"`
	SyncStatus   string     `gorm:"column:sync_status;type:varchar(50)"`
	ErrorMessage string     `gorm:"column:error_message;type:text"`

	// Standard timestamps
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName sets the table name
func (Mailbox) TableName() string {
	return "mailboxes"
}

func (m *Mailbox) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = utils.GenerateNanoIdWithPrefix("mbox", 16)
	}
	return nil
}
