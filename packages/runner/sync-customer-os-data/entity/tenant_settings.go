package entity

type TenantSettings struct {
	Tenant                 string `gorm:"column:tenant_name"`
	SalesforceClientId     string `gorm:"column:salesforce_client_id"`
	SalesforceClientSecret string `gorm:"column:salesforce_client_secret"`
	SalesforceRefreshToken string `gorm:"column:salesforce_refresh_token"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
