package entity

import "time"

type SlackSyncSettings struct {
	Tenant         string     `gorm:"type:varchar(255);not null"`
	ChannelId      string     `gorm:"varchar(50);not null"`
	ChannelName    string     `gorm:"varchar(1000)"`
	TeamId         string     `gorm:"varchar(50)"`
	SlackAccess    bool       `gorm:"column:slack_access;not null;default:true"`
	OrganizationId string     `gorm:"varchar(50)"`
	Enabled        bool       `gorm:"column:enabled;not null;default:false"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP - interval '30 day'"`
	SyncSince      time.Time  `gorm:"type:timestamp with time zone;not null;default:date_trunc('day', CURRENT_TIMESTAMP - interval '30 day')" `
	LookbackWindow int        `gorm:"default:7;not null"`
	LastSyncAt     *time.Time `gorm:"type:timestamp with time zone"`
}

func (SlackSyncSettings) TableName() string {
	return "slack_sync_settings"
}

func (s SlackSyncSettings) GetSyncStartDate() time.Time {
	if s.LastSyncAt == nil {
		return s.SyncSince
	}
	return *s.LastSyncAt
}
