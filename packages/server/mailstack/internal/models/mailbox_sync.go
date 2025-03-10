package models

import (
	"time"
)

// MailboxSyncState represents the synchronization state for a mailbox folder
type MailboxSyncState struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	MailboxID  string    `gorm:"column:mailbox_id;type:varchar(50);index;not null"`
	FolderName string    `gorm:"column:folder_name;type:varchar(100);index;not null"`
	LastUID    uint32    `gorm:"column:last_uid;not null"`
	LastSync   time.Time `gorm:"column:last_sync;type:timestamp;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;default:current_timestamp"`
}

func (MailboxSyncState) TableName() string {
	return "mailbox_sync_states"
}
