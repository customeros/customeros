package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type EmailMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`

	SentAt *time.Time `gorm:"column:sent_at;type:timestamp;DEFAULT:current_timestamp"`
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

func (e EmailMessage) MarshalJSON() ([]byte, error) {
	type Alias EmailMessage
	return json.Marshal(&struct {
		*Alias
		To  []string
		Cc  []string
		Bcc []string
	}{
		Alias: (*Alias)(&e),
		To:    utils.StringToSlice(e.ToString),
		Cc:    utils.StringToSlice(e.CcString),
		Bcc:   utils.StringToSlice(e.BccString),
	})
}

func (e *EmailMessage) UnmarshalJSON(data []byte) error {
	type Alias EmailMessage

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	e.ToString = utils.SliceToString(aux.To)
	e.CcString = utils.SliceToString(aux.Cc)
	e.BccString = utils.SliceToString(aux.Bcc)

	return nil
}

func (EmailMessage) TableName() string {
	return "email_message"
}
