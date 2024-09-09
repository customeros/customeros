package entity

import "time"

type EmailValidationRequestBulkStatus string

const (
	EmailValidationRequestBulkStatusProcessing EmailValidationRequestBulkStatus = "processing"
	EmailValidationRequestBulkStatusCompleted  EmailValidationRequestBulkStatus = "completed"
)

type EmailValidationRequestBulk struct {
	RequestID           string                           `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"requestId"`
	Tenant              string                           `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenantId"`
	TotalEmails         int                              `gorm:"column:total_emails;type:int;NOT NULL" json:"totalEmails"`
	DeliverableEmails   int                              `gorm:"column:undeliverable_emails;type:int;DEFAULT:0" json:"deliverableEmails"`
	UndeliverableEmails int                              `gorm:"column:undeliverable_emails;type:int;DEFAULT:0" json:"undeliverableEmails"`
	Status              EmailValidationRequestBulkStatus `gorm:"column:status;type:varchar(50);NOT NULL" json:"status"`
	CreatedAt           time.Time                        `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt           time.Time                        `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	FileName            string                           `gorm:"column:file_name;type:varchar(255);NOT NULL" json:"fileName"`
	Priority            int                              `gorm:"column:priority;type:int;DEFAULT:0" json:"priority"`
	VerifyCatchAll      bool                             `gorm:"column:verify_catch_all;type:boolean;DEFAULT:false" json:"verifyCatchAll"`
}

func (EmailValidationRequestBulk) TableName() string {
	return "email_validation_request_bulk"
}
