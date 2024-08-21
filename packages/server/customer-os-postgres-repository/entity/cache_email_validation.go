package entity

import "time"

type CacheEmailValidation struct {
	ID              string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Email           string    `gorm:"column:email;type:varchar(255);NOT NULL;index:idx_cache_email_validation_email,unique" json:"email"`
	NormalizedEmail string    `gorm:"column:normalized_email;type:varchar(255)" json:"normalizedEmail"`
	Username        string    `gorm:"column:username;type:varchar(255)" json:"username"`
	Domain          string    `gorm:"column:domain;type:varchar(255)" json:"domain"`
	IsDeliverable   bool      `gorm:"column:is_deliverable;type:boolean" json:"isDeliverable"`
	IsMailboxFull   bool      `gorm:"column:is_mailbox_full;type:boolean" json:"isMailboxFull"`
	IsRoleAccount   bool      `gorm:"column:is_role_account;type:boolean" json:"isRoleAccount"`
	IsFreeAccount   bool      `gorm:"column:is_free_account;type:boolean" json:"isFreeAccount"`
	SmtpSuccess     bool      `gorm:"column:smtp_success;type:boolean" json:"smtpSuccess"`
	ResponseCode    string    `gorm:"column:response_code;type:varchar(255)" json:"responseCode"`
	ErrorCode       string    `gorm:"column:error_code;type:varchar(255)" json:"errorCode"`
	Description     string    `gorm:"column:description;type:text" json:"description"`
	RetryValidation bool      `gorm:"column:retry_validation;type:boolean" json:"retryValidation"`
	SmtpResponse    string    `gorm:"column:smtp_response;type:text" json:"smtpResponse"`
}

func (CacheEmailValidation) TableName() string {
	return "cache_email_validation"
}
