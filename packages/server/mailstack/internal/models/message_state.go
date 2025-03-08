package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageState struct {
	ID            uint      `gorm:"primaryKey"`
	MailboxID     string    `gorm:"type:varchar(50);not null;index:idx_mailbox_folder,unique"`
	FolderName    string    `gorm:"type:varchar(255);not null;index:idx_mailbox_folder,unique"`
	LastSeenUID   uint32    `gorm:"not null;default:0"`
	LastCheckedAt time.Time `gorm:"not null;default:now()"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
