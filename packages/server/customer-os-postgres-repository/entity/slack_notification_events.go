package entity

import (
	"github.com/google/uuid"
	"time"
)

type SlackNotificationEvents struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant       string     `gorm:"index:idx_tenant;size:255;not null"`
	ChannelID    string     `gorm:"type:varchar(55)" json:"channelId"`
	Domain       string     `gorm:"index;type:varchar(55)" json:"domain"`
	LastNotified *time.Time `gorm:"type:timestamp;default:current_timestamp" json:"lastNotified"`
}

func (SlackNotificationEvents) TableName() string {
	return "slack_notification_events"
}
