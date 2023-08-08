package entity

type TenantSettings struct {
	ID         string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	TenantName string `gorm:"column:tenant_name;type:varchar(255);NOT NULL" binding:"required"`

	SlackApiToken       *string `gorm:"column:slack_api_token;type:varchar(255);" binding:"required"`
	SlackChannelFilter  *string `gorm:"column:slack_channel_filter;type:varchar(255);" binding:"required"`
	SlackLookbackWindow *string `gorm:"column:slack_lookback_window;type:varchar(255);" binding:"required"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
