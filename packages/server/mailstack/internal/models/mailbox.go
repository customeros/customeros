package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Mailbox struct {
	ID       string             `gorm:"column:id;type:varchar(50);primaryKey" json:"id"`
	Provider enum.EmailProvider `gorm:"column:provider;type:varchar(50);index;not null" json:"provider"`
	// IMAP Configuration
	ImapServer   string `gorm:"column:imap_server;type:varchar(255);not null" json:"imapServer"`
	ImapPort     int    `gorm:"column:imap_port;not null" json:"imapPort"`
	ImapUsername string `gorm:"column:imap_username;type:varchar(255);not null" json:"imapUsername"`
	ImapPassword string `gorm:"column:imap_password;type:varchar(255);not null" json:"imapPassword"`
	ImapTLS      bool   `gorm:"column:imap_tls;not null;default:true" json:"imapTls"`
	// SMTP Configuration
	SmtpServer   string `gorm:"column:smtp_server;type:varchar(255);not null" json:"smtpServer"`
	SmtpPort     int    `gorm:"column:smtp_port;not null" json:"smtpPort"`
	SmtpUsername string `gorm:"column:smtp_username;type:varchar(255);not null" json:"smtpUsername"`
	SmtpPassword string `gorm:"column:smtp_password;type:varchar(255);not null" json:"smtpPassword"`
	SmtpTLS      bool   `gorm:"column:smtp_tls;not null;default:true" json:"smtpTls"`
	// Other Configuration
	Folders      pq.StringArray `gorm:"column:folders;type:text[];not null" json:"folders"`
	DisplayName  string         `gorm:"column:display_name;type:varchar(255)" json:"displayName"`
	EmailAddress string         `gorm:"column:email_address;type:varchar(255);index" json:"emailAddress"`
	// Status Information
	LastSynced   *time.Time `gorm:"column:last_synced;type:timestamp" json:"lastSynced"`
	SyncStatus   string     `gorm:"column:sync_status;type:varchar(50)" json:"syncStatus"`
	ErrorMessage string     `gorm:"column:error_message;type:text" json:"errorMessage"`
	// Standard timestamps
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:current_timestamp" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
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
