package entity

import "time"

type TenantWebhook struct {
	// tenant, event, webhook, api key
	ID              uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	TenantName      string    `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName" binding:"required"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	WebhookUrl      string    `gorm:"column:webhook_url;type:varchar(255);NOT NULL" json:"webhook" binding:"required"`
	ApiKey          string    `gorm:"column:api_key;type:varchar(255);NOT NULL" json:"apiKey" binding:"required"`
	Event           string    `gorm:"column:event;type:varchar(255);NOT NULL" json:"event" binding:"required"`
	AuthHeaderName  string    `gorm:"column:auth_header_name;type:varchar(255)" json:"authHeaderName"`
	AuthHeaderValue string    `gorm:"column:auth_header_value;type:varchar(255)" json:"authHeaderValue"`
	// data for notifying user if webhook fails
	UserId        string `gorm:"column:user_id;type:varchar(255)" json:"userId"`
	UserFirstName string `gorm:"column:user_first_name;type:varchar(255)" json:"userFirstName"`
	UserLastName  string `gorm:"column:user_last_name;type:varchar(255)" json:"userLastName"`
	UserEmail     string `gorm:"column:user_email;type:varchar(255)" json:"userEmail"`
}

func (TenantWebhook) TableName() string {
	return "tenant_webhooks"
}

func (TenantWebhook) UniqueIndex() [][]string {
	return [][]string{
		{"TenantName", "Event"},
	}
}
