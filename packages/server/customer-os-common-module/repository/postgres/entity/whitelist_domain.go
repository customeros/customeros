package entity

import "time"

type WhitelistDomain struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"tenant" binding:"required"`
	Name      string    `gorm:"column:name;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"name" binding:"required"`
	Domain    string    `gorm:"column:domain;type:varchar(100);NOT NULL;index:name_domain_idx,unique" json:"domain" binding:"required"`
	Allowed   bool      `gorm:"column:allowed;" json:"allowed" binding:"required"`
}

func (WhitelistDomain) TableName() string {
	return "whitelist_domains"
}
