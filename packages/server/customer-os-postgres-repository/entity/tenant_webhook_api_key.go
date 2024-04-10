package entity

import "time"

type TenantWebhookApiKey struct {
	ID         uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	TenantName string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName" binding:"required"`
	Key        string    `gorm:"column:key;type:varchar(255);NOT NULL;index:idx_key,unique" json:"key" binding:"required"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
}

func (TenantWebhookApiKey) TableName() string {
	return "tenant_webhook_api_keys"
}
