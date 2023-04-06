package entity

type TenantSettings struct {
	ID         string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	TenantName string `gorm:"column:tenant_name;type:varchar(255);NOT NULL" binding:"required"`

	HubspotPrivateAppKey *string `gorm:"column:hubspot_private_app_key;type:varchar(255);" binding:"required"`

	ZendeskAPIKey     *string `gorm:"column:zendesk_api_key;type:varchar(255);" binding:"required"`
	ZendeskSubdomain  *string `gorm:"column:zendesk_subdomain;type:varchar(255);" binding:"required"`
	ZendeskAdminEmail *string `gorm:"column:zendesk_admin_email;type:varchar(255);" binding:"required"`

	SmartSheetId          *string `gorm:"column:smart_sheet_id;type:varchar(255);" binding:"required"`
	SmartSheetAccessToken *string `gorm:"column:smart_sheet_access_token;type:varchar(255);"  binding:"required"`

	JiraAPIToken *string `gorm:"column:jira_api_token;type:varchar(255);" binding:"required"`
	JiraDomain   *string `gorm:"column:jira_domain;type:varchar(255);" binding:"required"`
	JiraEmail    *string `gorm:"column:jira_email;type:varchar(255);" binding:"required"`

	TrelloAPIToken *string `gorm:"column:trello_api_token;type:varchar(255);" binding:"required"`
	TrelloAPIKey   *string `gorm:"column:trello_api_key;type:varchar(255);"  binding:"required"`

	AhaAPIUrl *string `gorm:"column:aha_api_url;type:varchar(255);"  binding:"required"`
	AhaAPIKey *string `gorm:"column:aha_api_key;type:varchar(255);"  binding:"required"`

	AirtablePersonalAccessToken *string `gorm:"column:airtable_personal_access_token;type:varchar(255);"  binding:"required"`

	AmplitudeAPIKey    *string `gorm:"column:amplitude_api_key;type:varchar(255);"  binding:"required"`
	AmplitudeSecretKey *string `gorm:"column:amplitude_secret_key;type:varchar(255);" binding:"required"`

	AsanaAccessToken *string `gorm:"column:asana_access_token;type:varchar(255);"  binding:"required"`

	BatonAPIKey *string `gorm:"column:baton_api_key;type:varchar(255);"  binding:"required"`

	BabelforceRegionEnvironment *string `gorm:"column:babelforce_region_environment;type:varchar(255);"  binding:"required"`
	BabelforceAccessKeyId       *string `gorm:"column:babelforce_access_key_id;type:varchar(255);"  binding:"required"`
	BabelforceAccessToken       *string `gorm:"column:babelforce_access_token;type:varchar(255);" binding:"required"`

	BigQueryServiceAccountKey *string `gorm:"column:bigquery_service_account_key;type:varchar(255);" binding:"required"`

	BraintreePublicKey   *string `gorm:"column:braintree_public_key;type:varchar(255);" binding:"required"`
	BraintreePrivateKey  *string `gorm:"column:braintree_private_key;type:varchar(255);"  binding:"required"`
	BraintreeEnvironment *string `gorm:"column:braintree_environment;type:varchar(255);"  binding:"required"`
	BraintreeMerchantId  *string `gorm:"column:braintree_merchant_id;type:varchar(255);"  binding:"required"`

	CallRailAccount  *string `gorm:"column:callrail_account;type:varchar(255);" binding:"required"`
	CallRailApiToken *string `gorm:"column:callrail_api_token;type:varchar(255);"  binding:"required"`

	ChargebeeApiKey         *string `gorm:"column:chargebee_api_key;type:varchar(255);" binding:"required"`
	ChargebeeProductCatalog *string `gorm:"column:chargebee_product_catalog;type:varchar(255);" binding:"required"`

	ChargifyApiKey *string `gorm:"column:chargify_api_key;type:varchar(255);" binding:"required"`
	ChargifyDomain *string `gorm:"column:chargify_domain;type:varchar(255);" binding:"required"`

	ClickUpApiKey *string `gorm:"column:clickup_api_key;type:varchar(255);" binding:"required"`

	CloseComApiKey *string `gorm:"column:closecom_api_key;type:varchar(255);" binding:"required"`

	CodaAuthToken  *string `gorm:"column:coda_auth_token;type:varchar(255);" binding:"required"`
	CodaDocumentId *string `gorm:"column:coda_document_id;type:varchar(255);" binding:"required"`

	ConfluenceApiToken   *string `gorm:"column:confluence_api_token;type:varchar(255);" binding:"required"`
	ConfluenceDomain     *string `gorm:"column:confluence_domain;type:varchar(255);" binding:"required"`
	ConfluenceLoginEmail *string `gorm:"column:confluence_login_email;type:varchar(255);" binding:"required"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
