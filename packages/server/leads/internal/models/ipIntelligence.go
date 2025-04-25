package models

import "time"

type IPIntelligence struct {
	ID           string     `gorm:"column:id;primaryKey;varchar(50)"`
	IPAddress    string     `gorm:"column:ip_address:varchar(50);not null"`
	Domain       string     `gorm:"column:domain:varchar(255);not null"`
	DomainSource string     `gorm:"column:domainSource:varchar(255);not null"`
	EmailAddress string     `gorm:"column:email_address:varchar(255);not null"`
	IsMobile     bool       `gorm:"column:is_mobile;type:bool;default:false;not null"`
	City         string     `gorm:"column:city:varchar(155);not null"`
	Region       string     `gorm:"column:region:varchar(155);not null"`
	CountryCode  string     `gorm:"column:country_code:varchar(10);not null"`
	HasThreat    bool       `gorm:"column:has_threat:bool;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (IPIntelligence) TableName() string {
	return "ip_intelligence"
}
