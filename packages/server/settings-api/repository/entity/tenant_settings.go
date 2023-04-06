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

	JiraAPIToken *string `gorm:"column:jira_api_token;type:varchar(255);" json:"jiraAPIToken" binding:"required"`
	JiraDomain   *string `gorm:"column:jira_domain;type:varchar(255);" json:"jiraDomain" binding:"required"`
	JiraEmail    *string `gorm:"column:jira_email;type:varchar(255);" json:"jiraEmail" binding:"required"`

	TrelloAPIToken *string `gorm:"column:trello_api_token;type:varchar(255);" json:"trelloAPIToken" binding:"required"`
	TrelloAPIKey   *string `gorm:"column:trello_api_key;type:varchar(255);" json:"trelloAPIKey" binding:"required"`

	AhaAPIUrl *string `gorm:"column:aha_api_url;type:varchar(255);" json:"ahaAPIUrl" binding:"required"`
	AhaAPIKey *string `gorm:"column:aha_api_key;type:varchar(255);" json:"ahaAPIKey" binding:"required"`

	AirtablePersonalAccessToken *string `gorm:"column:airtable_personal_access_token;type:varchar(255);" json:"airtablePersonalAccessToken" binding:"required"`

	AmplitudeAPIKey    *string `gorm:"column:amplitude_api_key;type:varchar(255);" json:"amplitudeAPIKey" binding:"required"`
	AmplitudeSecretKey *string `gorm:"column:amplitude_secret_key;type:varchar(255);" json:"amplitudeSecretKey" binding:"required"`

	BatonAPIKey *string `gorm:"column:baton_api_key;type:varchar(255);" json:"batonAPIKey" binding:"required"`

	BabelforceRegionEnvironment *string `gorm:"column:babelforce_region_environment;type:varchar(255);" json:"babelforceRegionEnvironment" binding:"required"`
	BabelforceAccessKeyId       *string `gorm:"column:babelforce_access_key_id;type:varchar(255);" json:"babelforceAccessKeyId" binding:"required"`
	BabelforceAccessToken       *string `gorm:"column:babelforce_access_token;type:varchar(255);" json:"babelforceAccessToken" binding:"required"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
