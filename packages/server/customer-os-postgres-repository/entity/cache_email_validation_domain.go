package entity

import "time"

type CacheEmailValidationDomain struct {
	ID                  string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Domain              string    `gorm:"column:domain;type:varchar(255);NOT NULL;index:idx_cache_email_validation_domain_domain,unique" json:"domain"`
	IsCatchAll          bool      `gorm:"column:is_catch_all;type:boolean" json:"isCatchAll"`
	IsFirewalled        bool      `gorm:"column:is_firewalled;type:boolean" json:"isFirewalled"`
	Provider            string    `gorm:"column:provider;type:varchar(255)" json:"provider"`
	Firewall            string    `gorm:"column:firewall;type:varchar(255)" json:"firewall"`
	HasMXRecord         bool      `gorm:"column:has_mx_record;type:boolean" json:"hasMXRecord"`
	HasSPFRecord        bool      `gorm:"column:has_spf_record;type:boolean" json:"hasSPFRecord"`
	Error               string    `gorm:"column:error;type:varchar(255)" json:"error"`
	Data                string    `gorm:"column:data;type:text" json:"data"`
	CanConnectSMTP      bool      `gorm:"column:can_connect_smtp;type:boolean" json:"canConnectSMTP"`
	TLSRequired         bool      `gorm:"column:tls_required;type:boolean" json:"tlsRequired"`
	ResponseCode        string    `gorm:"column:response_code;type:varchar(255)" json:"responseCode"`
	ErrorCode           string    `gorm:"column:error_code;type:varchar(255)" json:"errorCode"`
	Description         string    `gorm:"column:description;type:text" json:"description"`
	HealthIsGreylisted  bool      `gorm:"column:health_is_greylisted;type:boolean" json:"healthIsGreylisted"`
	HealthIsBlacklisted bool      `gorm:"column:health_is_blacklisted;type:boolean" json:"healthIsBlacklisted"`
	HealthServerIP      string    `gorm:"column:health_server_ip;type:varchar(255)" json:"healthServerIP"`
	HealthFromEmail     string    `gorm:"column:health_from_email;type:varchar(255)" json:"healthFromEmail"`
	HealthRetryAfter    int       `gorm:"column:health_retry_after;type:integer" json:"healthRetryAfter"`
	IsPrimaryDomain     *bool     `gorm:"column:is_primary_domain;type:boolean" json:"isPrimaryDomain"`
	PrimaryDomain       string    `gorm:"column:primary_domain;type:varchar(255)" json:"primaryDomain"`
}

func (CacheEmailValidationDomain) TableName() string {
	return "cache_email_validation_domain"
}
