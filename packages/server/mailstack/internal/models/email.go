package models

import (
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Email represents a raw email message stored in the database
type Email struct {
	ID         string             `gorm:"column:id;type:varchar(50);primaryKey"`
	MailboxID  string             `gorm:"column:mailbox_id;type:varchar(50);index;not null"`
	Provider   enum.EmailProvider `gorm:"column:provider;type:varchar(50);index;not null"`
	Folder     string             `gorm:"column:folder;type:varchar(100);index;not null"`
	ImapUID    uint32             `gorm:"column:imap_uid;index"`
	MessageID  string             `gorm:"column:message_id;uniqueIndex;type:varchar(255);not null"`
	ThreadID   string             `gorm:"column:thread_id;type:varchar(255);index"`
	InReplyTo  string             `gorm:"column:in_reply_to;type:varchar(255);index"`
	References pq.StringArray     `gorm:"column:references;type:text[]"`

	// Core email metadata
	Subject      string         `gorm:"column:subject;type:varchar(1000)"`
	FromAddress  string         `gorm:"column:from_address;type:varchar(255);index"`
	FromName     string         `gorm:"column:from_name;type:varchar(255)"`
	ReplyTo      string         `gorm:"column:reply_to;type:varchar(255);index"`
	ToAddresses  pq.StringArray `gorm:"column:to_addresses;type:text[]"`
	CcAddresses  pq.StringArray `gorm:"column:cc_addresses;type:text[]"`
	BccAddresses pq.StringArray `gorm:"column:bcc_addresses;type:text[]"`

	// Time information
	SentAt     *time.Time `gorm:"column:sent_at;type:timestamp;index"`
	ReceivedAt *time.Time `gorm:"column:received_at;type:timestamp;index"`

	// Content
	BodyText      string `gorm:"column:body_text;type:text"`
	BodyHTML      string `gorm:"column:body_html;type:text"`
	HasAttachment bool   `gorm:"column:has_attachment;default:false"`

	// Extensions and provider-specific data
	GmailLabels       pq.StringArray `gorm:"column:gmail_labels;type:text[]"`
	OutlookCategories pq.StringArray `gorm:"column:outlook_categories;type:text[]"`
	MailstackFlags    pq.StringArray `gorm:"column:mailstack_flags;type:text[]"`

	// Raw data
	RawHeaders    JSONMap `gorm:"column:raw_headers;type:jsonb"`
	Envelope      JSONMap `gorm:"column:envelope;type:jsonb"`
	BodyStructure JSONMap `gorm:"column:body_structure;type:jsonb"`

	// Classification
	Classification       enum.EmailClassification `gorm:"column:classification;type:varchar(50);index"`
	ClassificationReason string                   `gorm:"column:classification_reason;type:text"`

	// Standard timestamps
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

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

// EmailHeaders represents specific email header information needed for processing
type EmailHeaders struct {
	AutoSubmitted      bool
	ContentDescription string
	DeliveryStatus     bool
	ListUnsubscribe    bool
	Precedence         string
	ReturnPath         string
	ReturnPathExists   bool
	XAutoreply         string
	XAutoresponse      string
	XLoop              bool
	XFailedRecipients  []string
	ReplyTo            string
	ReplyToExists      bool
	Sender             string
	ForwardedFor       string
	DKIM               []string
	SPF                string
	DMARC              string
}

func (e *Email) Headers() (*EmailHeaders, error) {
	headers := &EmailHeaders{}

	if e.RawHeaders == nil {
		return headers, nil
	}

	// Helper function to get header value as string
	getString := func(key string) string {
		if values, ok := e.RawHeaders[key].([]string); ok && len(values) > 0 {
			return values[0]
		}
		if value, ok := e.RawHeaders[key].(string); ok {
			return value
		}
		return ""
	}

	// Helper function to check if a header exists
	headerExists := func(key string) bool {
		_, exists := e.RawHeaders[key]
		return exists
	}

	// Helper to get string array
	getStringArray := func(key string) []string {
		if values, ok := e.RawHeaders[key].([]string); ok {
			return values
		}
		if value, ok := e.RawHeaders[key].(string); ok {
			return []string{value}
		}
		return nil
	}

	// Process boolean headers (presence/absence or specific values)
	autoSubmitted := getString("Auto-Submitted")
	headers.AutoSubmitted = autoSubmitted != "" && autoSubmitted != "no"

	// Content-Description
	headers.ContentDescription = getString("Content-Description")

	// Delivery-Status
	headers.DeliveryStatus = headerExists("Delivery-Status") ||
		headerExists("X-Failed-Recipients") ||
		strings.Contains(e.Subject, "Delivery Status Notification") ||
		strings.Contains(e.Subject, "Mail Delivery Failure")

	// List-Unsubscribe
	headers.ListUnsubscribe = headerExists("List-Unsubscribe")

	// Precedence
	headers.Precedence = getString("Precedence")

	// Return-Path
	returnPath := getString("Return-Path")
	headers.ReturnPath = returnPath
	headers.ReturnPathExists = headerExists("Return-Path")

	// Auto-reply headers
	headers.XAutoreply = getString("X-Autoreply")
	headers.XAutoresponse = getString("X-Autoresponse")

	// X-Loop
	headers.XLoop = headerExists("X-Loop")

	// X-Failed-Recipients
	failedRecipientsStr := getString("X-Failed-Recipients")
	if failedRecipientsStr != "" {
		recipients := strings.Split(failedRecipientsStr, ",")
		for i, recipient := range recipients {
			recipients[i] = strings.TrimSpace(recipient)
		}
		headers.XFailedRecipients = recipients
	}

	// Reply-To
	headers.ReplyTo = getString("Reply-To")
	headers.ReplyToExists = headerExists("Reply-To")

	// Sender
	headers.Sender = getString("Sender")

	// Forwarded-For (could be in different formats)
	headers.ForwardedFor = getString("X-Forwarded-For")
	if headers.ForwardedFor == "" {
		headers.ForwardedFor = getString("Forwarded-For")
	}

	// Security headers
	headers.DKIM = getStringArray("DKIM-Signature")
	headers.SPF = getString("Received-SPF")
	headers.DMARC = getString("DMARC-Result")

	return headers, nil
}
