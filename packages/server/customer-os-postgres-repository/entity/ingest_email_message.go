package postgres_entity

import (
	"time"
)

type IngestEmailMessage struct {
	Id        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`

	Error *string                 `gorm:"column:error;type:text"`
	State IngestEmailMessageState `gorm:"size:255;not null;"`

	Tenant   string `gorm:"size:255;not null;index:idx_ingest_raw_email"`
	Provider string `gorm:"size:255;not null;index:idx_ingest_raw_email"`
	Username string `gorm:"size:255;not null;index:idx_ingest_raw_email"`

	//Email message data
	SentAt time.Time `gorm:"column:sent_at;type:timestamp;"`

	From        string `gorm:"type:text"`
	To          string `gorm:"type:text"`
	Cc          string `gorm:"type:text"`
	Bcc         string `gorm:"type:text"`
	Subject     string
	TextContent string `gorm:"type:text"`
	HtmlContent string `gorm:"type:text"`

	//Values taken from providers
	MessageId          string `gorm:"size:255;not null;index:idx_ingest_raw_email"`
	ProviderMessageId  string `gorm:"size:255;not null;"`
	ProviderThreadId   string `gorm:"size:255;not null;"`
	ProviderInReplyTo  string `gorm:"type:text"`
	ProviderReferences string `gorm:"type:text"`

	Headers string `gorm:"type:text"`
}

func (IngestEmailMessage) TableName() string {
	return "ingest_email_message"
}

type IngestEmailMessageState string

const (
	IngestEmailMessageStateError       IngestEmailMessageState = "ERROR"
	IngestEmailMessageStatePending     IngestEmailMessageState = "PENDING"
	IngestEmailMessageStateSentToAgent IngestEmailMessageState = "SENT_TO_AGENT"
	IngestEmailMessageStateIngested    IngestEmailMessageState = "INGESTED"
	IngestEmailMessageStateFiltered    IngestEmailMessageState = "FILTERED"
)

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
