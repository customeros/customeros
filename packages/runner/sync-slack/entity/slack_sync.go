package entity

import "time"

type SlackSync struct {
	ChannelId  string    `gorm:"primary_key;varchar(50)" binding:"required"`
	Tenant     string    `gorm:"primary_key;type:varchar(255)" binding:"required"`
	LastSyncAt time.Time `gorm:"type:timestamp with time zone;not null" binding:"required"`
}

func (SlackSync) TableName() string {
	return "slack_sync"
}
