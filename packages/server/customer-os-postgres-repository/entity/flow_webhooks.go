package entity

import (
	"time"

	"github.com/pkg/errors"
)

type FlowWebhooks struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Tenant        string     `gorm:"column:tenant;type:varchar(255);not null;index" json:"tenant" binding:"required"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	WebhookPath   string     `gorm:"column:webhook_path;type:varchar(255);not null" json:"webhookPath" binding:"required"`
	Integration   string     `gorm:"column:integration;type:varchar(255);index" json:"integration"`
	RotationCount int        `gorm:"column:rotation_count;type:integer" json:"rotationCount"`
	Secret        string     `gorm:"column:secret;type:varchar(255);not null" json:"secret" binding:"required"`
	Enabled       bool       `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
}

func (FlowWebhooks) TableName() string {
	return "flow_webhooks"
}

func (FlowWebhooks) UniqueIndex() [][]string {
	return [][]string{
		{"tenant", "integration"},
	}
}

func (fw *FlowWebhooks) Validate() error {
	if fw.Tenant == "" {
		return errors.New("tenant is required")
	}
	if fw.Integration == "" {
		return errors.New("integration is required")
	}
	if fw.Secret == "" {
		return errors.New("secret is required")
	}
	return nil
}

func (fw *FlowWebhooks) IsActive() bool {
	return fw.Enabled && fw.Secret != ""
}
