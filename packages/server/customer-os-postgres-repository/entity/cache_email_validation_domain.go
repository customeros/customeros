package entity

import "time"

type CacheEmailValidationDomain struct {
	ID             string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Domain         string    `gorm:"column:domain;type:varchar(255);NOT NULL;index:idx_cache_email_validation_domain_domain,unique" json:"domain"`
	IsCatchAll     bool      `gorm:"column:is_catch_all;type:boolean" json:"isCatchAll"`
	IsFirewalled   bool      `gorm:"column:is_firewalled;type:boolean" json:"isFirewalled"`
	CanConnectSMTP bool      `gorm:"column:can_connect_smtp;type:boolean" json:"canConnectSMTP"`
	Provider       string    `gorm:"column:provider;type:varchar(255)" json:"provider"`
	Firewall       string    `gorm:"column:firewall;type:varchar(255)" json:"firewall"`
}

func (CacheEmailValidationDomain) TableName() string {
	return "cache_email_validation_domain"
}
