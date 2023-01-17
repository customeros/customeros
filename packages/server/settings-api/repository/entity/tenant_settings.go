package entity

type TenantSettings struct {
	ID         string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	TenantName string `gorm:"column:tenant_name;type:varchar(255);NOT NULL" json:"tenantName" binding:"required"`

	HubspotPrivateAppKey *string `gorm:"column:hubspot_private_app_key;type:varchar(255);" json:"hubspotPrivateAppKey" binding:"required"`

	ZendeskAPIKey     *string `gorm:"column:zendesk_api_key;type:varchar(255);" json:"zendeskAPIKey" binding:"required"`
	ZendeskSubdomain  *string `gorm:"column:zendesk_subdomain;type:varchar(255);" json:"zendeskSubdomain" binding:"required"`
	ZendeskAdminEmail *string `gorm:"column:zendesk_admin_email;type:varchar(255);" json:"zendeskAdminEmail" binding:"required"`

	SmartSheetId          *string `gorm:"column:smart_sheet_id;type:varchar(255);" json:"smartSheetId" binding:"required"`
	SmartSheetAccessToken *string `gorm:"column:smart_sheet_access_token;type:varchar(255);" json:"smartSheetAccessToken" binding:"required"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
