package entity

import "time"

type TenantWebhook struct {
	// tenant, event, webhook, api key
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	TenantName string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName" binding:"required"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	WebhookUrl string    `gorm:"column:webhook_url;type:varchar(255);NOT NULL" json:"webhook" binding:"required"`
	ApiKey     string    `gorm:"column:api_key;type:varchar(255);NOT NULL" json:"apiKey" binding:"required"`
	Event      string    `gorm:"column:event;type:varchar(255);NOT NULL" json:"event" binding:"required"`
	AuthHeader string    `gorm:"column:auth_header;type:varchar(255)" json:"authHeader"`
}

func (TenantWebhook) TableName() string {
	return "tenant_webhooks"
}

func (TenantWebhook) UniqueIndex() [][]string {
	return [][]string{
		{"TenantName", "Event"},
	}
}
