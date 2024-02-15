package entity

import "time"

type SlackChannel struct {
	ID             uint64    `gorm:"primary_key;autoIncrement:true"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp"`
	TenantName     string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL"`
	ChannelId      string    `gorm:"column:channel_id;type:varchar(255);NOT NULL"`
	OrganizationId *string   `gorm:"column:organization_id;type:varchar(255);"`
}

func (SlackChannel) TableName() string {
	return "slack_channel"
}

func (SlackChannel) UniqueIndex() [][]string {
	return [][]string{
		{"TenantName", "ChannelId"},
	}
}
