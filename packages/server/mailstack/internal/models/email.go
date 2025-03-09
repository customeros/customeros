package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Email represents a raw email message stored in the database
type Email struct {
	ID        string `gorm:"type:varchar(50);primaryKey"`
	MailboxID string `gorm:"type:varchar(50);index;not null"`
	Provider  string `gorm:"type:varchar(50);index;not null"`
	Folder    string `gorm:"type:varchar(100);index;not null"`
	ImapUID   uint32 `gorm:"index"`
	MessageID string `gorm:"type:varchar(255);index"`
	ThreadID  string `gorm:"type:varchar(255);index"`

	// Core email metadata
	Subject      string         `gorm:"type:varchar(1000)"`
	FromAddress  string         `gorm:"type:varchar(255);index"`
	FromName     string         `gorm:"type:varchar(255)"`
	ToAddresses  pq.StringArray `gorm:"type:text[]"`
	CcAddresses  pq.StringArray `gorm:"type:text[]"`
	BccAddresses pq.StringArray `gorm:"type:text[]"`

	// Time information
	ReceivedAt time.Time `gorm:"index"`
	SentAt     time.Time `gorm:"index"`

	// Content
	BodyText      string `gorm:"type:text"`
	BodyHTML      string `gorm:"type:text"`
	HasAttachment bool   `gorm:"default:false"`

	// Extensions and provider-specific data
	GmailLabels       pq.StringArray `gorm:"type:text[]"`
	OutlookCategories pq.StringArray `gorm:"type:text[]"`
	MailstackFlags    pq.StringArray `gorm:"type:text[]"`

	// Raw data
	RawHeaders    JSONMap `gorm:"type:jsonb"`
	Envelope      JSONMap `gorm:"type:jsonb"`
	BodyStructure JSONMap `gorm:"type:jsonb"`

	// Standard timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName overrides the table name for Email
func (Email) TableName() string {
	return "emails"
}

func (e *Email) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = utils.GenerateNanoIdWithPrefix("email", 24)
	}
	e.CreatedAt = utils.Now()
	return nil
}

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

// JSONMap represents a JSON object that can be stored in PostgreSQL
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
