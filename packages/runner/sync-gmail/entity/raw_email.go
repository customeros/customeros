package entity

import (
	"github.com/google/uuid"
	"time"
)

type RawEmail struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	SentAt    time.Time `gorm:"column:sent_at;type:timestamp;DEFAULT:current_timestamp" json:"sentAt"`

	ExternalSystem    string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	Tenant            string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	Username          string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	State             string `gorm:"size:255;not null;"`
	ProviderMessageId string `gorm:"size:255;not null;"`
	MessageId         string `gorm:"size:255;not null;index:idx_raw_email_external_system"`

	SentToEventStoreState  string  `gorm:"size:50;not null"`
	SentToEventStoreReason *string `gorm:"type:text"`
	SentToEventStoreError  *string `gorm:"type:text"`

	Data string `gorm:"type:text"`
}

func (RawEmail) TableName() string {
	return "raw_email"
}
