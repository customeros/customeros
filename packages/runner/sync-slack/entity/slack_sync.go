package entity

import "time"

type SlackSync struct {
	ChannelId  string    `gorm:"primary_key;varchar(50)" binding:"required"`
	TenantName string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" binding:"required"`
	LastSyncAt time.Time `gorm:"type:timestamp with time zone;not null" binding:"required"`
}

func (SlackSync) TableName() string {
	return "slack_sync"
}
