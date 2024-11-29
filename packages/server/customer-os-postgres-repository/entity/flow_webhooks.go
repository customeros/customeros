package entity

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
)

type FlowWebhooks struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantName    string    `gorm:"type:varchar(255);not null;index" json:"tenantName" binding:"required"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"type:timestamp" json:"updatedAt"`
	WebhookPath   string    `gorm:"type:varchar(255);not null" json:"webhookPath" binding:"required"`
	Integration   string    `gorm:"type:varchar(255);index" json:"integration"`
	RotationCount int       `gorm:"type:integer" json:"rotationCount"`
	Secret        string    `gorm:"type:varchar(255);not null" json:"secret" binding:"required"`
	Enabled       bool      `gorm:"type:boolean;default:true" json:"enabled"`
}

func (TenantWebhook) TableName() string {
	return "flow_webhooks"
}

func (FlowWebhooks) UniqueIndex() [][]string {
	return [][]string{
		{"tenant_name", "integration"},
	}
}

func (fw *FlowWebhooks) BeforeCreate() error {
	if fw.TenantName == "" {
		return errors.New("tenant name is required")
	}
	fw.CreatedAt = time.Now()
	fw.UpdatedAt = time.Now()
	if fw.RotationCount == 0 {
		fw.RotationCount = 1
	}
	fw.Enabled = true
	return nil
}

func (fw *FlowWebhooks) BeforeUpdate() error {
	fw.UpdatedAt = time.Now()
	return nil
}

// Rotate increments rotation count and generates new secret
func (fw *FlowWebhooks) Rotate() error {
	fw.RotationCount++
	fw.Secret = utils.GenerateSecret()
	return nil
}

func (fw *FlowWebhooks) Validate() error {
	if fw.TenantName == "" {
		return errors.New("tenant name is required")
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
