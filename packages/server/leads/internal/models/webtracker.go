package models

import (
	"time"
)

type WebTracker struct {
	ID                string     `gorm:"column:id;primaryKey;type:varchar(50);"`
	Tenant            string     `gorm:"column:tenant;type:varchar(255);index;not null"`
	Domain            string     `gorm:"column:domain;type:varchar(255);index;not null"`
	CNAMEHost         string     `gorm:"column:cname_host;type:varchar(255);not null"`
	CNAMETarget       string     `gorm:"column:cname_target;type:varchar(255);not null"`
	IsCNAMEConfigured bool       `gorm:"column:is_cname_configured;type:bool;default:false;not null"`
	IsProxyActive     bool       `gorm:"column:is_proxy_active;type:bool;default:false;not null"`
	LastEventAt       *time.Time `gorm:"column:last_event_at;"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         *time.Time `gorm:"column:updated_at;autoUpdateTime"`
	IsArchived        bool       `gorm:"column:is_archived;type:bool;default:false;not null"`
}

func (WebTracker) TableName() string {
	return "web_trackers"
}
