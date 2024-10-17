package entity

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
	"time"
)

type EmailMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`

	Status EmailMessageStatus `gorm:"size:50;not null;"`

	SentAt *time.Time `gorm:"column:sent_at;type:timestamp;"`
	Error  *string    `gorm:"column:error;type:text"`

	UniqueInternalIdentifier *string `gorm:"size:255;index:unique_internal_identifier"`

	Tenant       string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	ProducerId   string `gorm:"size:255;not null;"`
	ProducerType string `gorm:"size:255;not null;"`

	//Email message data
	From         string   `gorm:"size:255;"`
	FromProvider string   `gorm:"size:255;"`
	To           []string `gorm:"-"`
	ToString     string   `gorm:"type:text;column:to"`
	Cc           []string `gorm:"-"`
	CcString     string   `gorm:"type:text;column:cc"`
	Bcc          []string `gorm:"-"`
	BccString    string   `gorm:"type:text;column:bcc"`
	Subject      string
	Content      string

	//COS interaction event id
	ReplyTo *string

	//Values taken from providers
	ProviderMessageId  string `gorm:"size:255;not null;"`
	ProviderThreadId   string `gorm:"size:255;not null;"`
	ProviderInReplyTo  string `gorm:"size:255;not null;"`
	ProviderReferences string `gorm:"size:255;not null;"`
}

type EmailMessageStatus string

const (
	EmailMessageStatusScheduled EmailMessageStatus = "SCHEDULED"
	EmailMessageStatusSent      EmailMessageStatus = "SENT"
	EmailMessageStatusProcessed EmailMessageStatus = "PROCESSED"
	EmailMessageStatusError     EmailMessageStatus = "ERROR"
)

// BeforeSave hook for converting slices to strings
func (e *EmailMessage) BeforeSave(tx *gorm.DB) (err error) {
	e.ToString = utils.SliceToString(e.To)
	e.CcString = utils.SliceToString(e.Cc)
	e.BccString = utils.SliceToString(e.Bcc)
	return nil
}

// AfterFind hook for converting strings to slices
func (e *EmailMessage) AfterFind(tx *gorm.DB) (err error) {
	e.To = utils.StringToSlice(e.ToString)
	e.Cc = utils.StringToSlice(e.CcString)
	e.Bcc = utils.StringToSlice(e.BccString)
	return nil
}

func (EmailMessage) TableName() string {
	return "email_message"
}
