package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

// EmailAttachment represents an attachment to an email
type EmailAttachment struct {
	ID          string `gorm:"type:varchar(50);primaryKey"`
	EmailID     string `gorm:"type:varchar(50);index;not null"`
	Filename    string `gorm:"type:varchar(500)"`
	ContentType string `gorm:"type:varchar(255)"`
	ContentID   string `gorm:"type:varchar(255)"` // For inline attachments
	Size        int    `gorm:"default:0"`
	IsInline    bool   `gorm:"default:false"`

	// Storage options
	StorageKey string `gorm:"type:varchar(1000)"` // If stored in S3/blob storage

	// Standard timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the table name for EmailAttachment
func (EmailAttachment) TableName() string {
	return "email_attachments"
}

func (e *EmailAttachment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = utils.GenerateNanoIdWithPrefix("atch", 21)
	}
	e.CreatedAt = utils.Now()
	return nil
}
