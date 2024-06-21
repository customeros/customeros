package entity

import (
	"github.com/google/uuid"
	"time"
)

type RawEmail struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	SentAt    time.Time `gorm:"column:sent_at;type:timestamp;DEFAULT:current_timestamp" json:"sentAt"`

	ExternalSystem    string           `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	Tenant            string           `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	Username          string           `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	State             EmailImportState `gorm:"size:255;not null;"`
	ProviderMessageId string           `gorm:"size:255;not null;"`
	MessageId         string           `gorm:"size:255;not null;index:idx_raw_email_external_system"`

	SentToEventStoreState  string  `gorm:"size:50;not null"`
	SentToEventStoreReason *string `gorm:"type:text"`
	SentToEventStoreError  *string `gorm:"type:text"`

	Data string `gorm:"type:text"`
}

func (RawEmail) TableName() string {
	return "raw_email"
}

type EmailRawData struct {
	ProviderMessageId string            `json:"ProviderMessageId"`
	MessageId         string            `json:"MessageId"`
	Sent              time.Time         `json:"Sent"`
	Subject           string            `json:"Subject"`
	From              string            `json:"From"`
	To                string            `json:"To"`
	Cc                string            `json:"Cc"`
	Bcc               string            `json:"Bcc"`
	Html              string            `json:"Html"`
	Text              string            `json:"Text"`
	ThreadId          string            `json:"ThreadId"`
	InReplyTo         string            `json:"InReplyTo"`
	Reference         string            `json:"Reference"`
	Headers           map[string]string `json:"Headers"`
}
