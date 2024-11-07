package entity

import "time"

type CacheEmailValidation struct {
	ID                  string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Email               string    `gorm:"column:email;type:varchar(255);NOT NULL;index:idx_cache_email_validation_email,unique" json:"email"`
	NormalizedEmail     string    `gorm:"column:normalized_email;type:varchar(255)" json:"normalizedEmail"`
	Username            string    `gorm:"column:username;type:varchar(255)" json:"username"`
	Domain              string    `gorm:"column:domain;type:varchar(255)" json:"domain"`
	IsMailboxFull       bool      `gorm:"column:is_mailbox_full;type:boolean" json:"isMailboxFull"`
	IsRoleAccount       bool      `gorm:"column:is_role_account;type:boolean" json:"isRoleAccount"`
	IsSystemGenerated   bool      `gorm:"column:is_system_generated;type:boolean" json:"isSystemGenerated"`
	IsFreeAccount       bool      `gorm:"column:is_free_account;type:boolean" json:"isFreeAccount"`
	SmtpSuccess         bool      `gorm:"column:smtp_success;type:boolean" json:"smtpSuccess"`
	ResponseCode        string    `gorm:"column:response_code;type:varchar(255)" json:"responseCode"`
	ErrorCode           string    `gorm:"column:error_code;type:varchar(255)" json:"errorCode"`
	Description         string    `gorm:"column:description;type:text" json:"description"`
	TLSRequired         bool      `gorm:"column:tls_required;type:boolean" json:"tlsRequired"`
	RetryValidation     bool      `gorm:"column:retry_validation;type:boolean" json:"retryValidation"`
	Deliverable         string    `gorm:"column:deliverable;type:varchar(16)" json:"deliverable"`
	HealthIsGreylisted  bool      `gorm:"column:health_is_greylisted;type:boolean" json:"healthIsGreylisted"`
	HealthIsBlacklisted bool      `gorm:"column:health_is_blacklisted;type:boolean" json:"healthIsBlacklisted"`
	HealthServerIP      string    `gorm:"column:health_server_ip;type:varchar(255)" json:"healthServerIP"`
	HealthFromEmail     string    `gorm:"column:health_from_email;type:varchar(255)" json:"healthFromEmail"`
	HealthRetryAfter    int       `gorm:"column:health_retry_after;type:integer" json:"healthRetryAfter"`
	AlternateEmail      string    `gorm:"column:alternate_email;type:varchar(255)" json:"alternateEmail"`
	Error               string    `gorm:"column:error;type:varchar(255)" json:"error"`
	Data                string    `gorm:"column:data;type:text" json:"data"`
}

func (CacheEmailValidation) TableName() string {
	return "cache_email_validation"
}
