package entity

type TenantSettings struct {
	ID         string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	TenantName string `gorm:"column:tenant_name;type:varchar(255);NOT NULL" binding:"required"`

	HubspotPrivateAppKey *string `gorm:"column:hubspot_private_app_key;type:varchar(255);" binding:"required"`

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

	CourierApiKey *string `gorm:"column:courier_api_key;type:varchar(255);" binding:"required"`

	CustomerIoApiKey *string `gorm:"column:customerio_api_key;type:varchar(255);" binding:"required"`

	DatadogApiKey         *string `gorm:"column:datadog_api_key;type:varchar(255);" binding:"required"`
	DatadogApplicationKey *string `gorm:"column:datadog_application_key;type:varchar(255);" binding:"required"`

	DelightedApiKey *string `gorm:"column:delighted_api_key;type:varchar(255);" binding:"required"`

	DixaApiToken *string `gorm:"column:dixa_api_token;type:varchar(255);" binding:"required"`

	DriftApiToken *string `gorm:"column:drift_api_token;type:varchar(255);" binding:"required"`

	EmailOctopusApiKey *string `gorm:"column:emailoctopus_api_key;type:varchar(255);" binding:"required"`

	FacebookMarketingAccessToken *string `gorm:"column:facebook_marketing_access_token;type:varchar(255);"  binding:"required"`

	FastbillApiKey    *string `gorm:"column:fastbill_api_key;type:varchar(255);" binding:"required"`
	FastbillProjectId *string `gorm:"column:fastbill_project_id;type:varchar(255);" binding:"required"`

	FlexportApiKey *string `gorm:"column:flexport_api_key;type:varchar(255);" binding:"required"`

	FreshcallerApiKey *string `gorm:"column:freshcaller_api_key;type:varchar(255);" binding:"required"`

	FreshdeskApiKey *string `gorm:"column:freshdesk_api_key;type:varchar(255);" binding:"required"`
	FreshdeskDomain *string `gorm:"column:freshdesk_domain;type:varchar(255);" binding:"required"`

	FreshsalesApiKey *string `gorm:"column:freshsales_api_key;type:varchar(255);" binding:"required"`
	FreshsalesDomain *string `gorm:"column:freshsales_domain;type:varchar(255);" binding:"required"`

	FreshserviceApiKey *string `gorm:"column:freshservice_api_key;type:varchar(255);" binding:"required"`
	FreshserviceDomain *string `gorm:"column:freshservice_domain;type:varchar(255);" binding:"required"`

	GenesysRegion       *string `gorm:"column:genesys_region;type:varchar(255);" binding:"required"`
	GenesysClientId     *string `gorm:"column:genesys_client_id;type:varchar(255);" binding:"required"`
	GenesysClientSecret *string `gorm:"column:genesys_client_secret;type:varchar(255);" binding:"required"`

	GitHubAccessToken *string `gorm:"column:github_access_token;type:varchar(255);" binding:"required"`

	GitLabAccessToken *string `gorm:"column:gitlab_access_token;type:varchar(255);" binding:"required"`

	GoCardlessAccessToken *string `gorm:"column:gocardless_access_token;type:varchar(255);" binding:"required"`
	GoCardlessEnvironment *string `gorm:"column:gocardless_environment;type:varchar(255);" binding:"required"`
	GoCardlessVersion     *string `gorm:"column:gocardless_version;type:varchar(255);" binding:"required"`

	GongApiKey *string `gorm:"column:gong_api_key;type:varchar(255);" binding:"required"`

	HarvestAccountId   *string `gorm:"column:harvest_account_id;type:varchar(255);" binding:"required"`
	HarvestAccessToken *string `gorm:"column:harvest_access_token;type:varchar(255);" binding:"required"`

	InsightlyApiToken *string `gorm:"column:insightly_api_token;type:varchar(255);" binding:"required"`

	InstagramAccessToken *string `gorm:"column:instagram_access_token;type:varchar(255);" binding:"required"`

	InstatusApiKey *string `gorm:"column:instatus_api_key;type:varchar(255);" binding:"required"`

	IntercomAccessToken *string `gorm:"column:intercom_access_token;type:varchar(255);" binding:"required"`

	KlaviyoApiKey *string `gorm:"column:klaviyo_api_key;type:varchar(255);" binding:"required"`

	KustomerApiToken *string `gorm:"column:kustomer_api_token;type:varchar(255);" binding:"required"`

	LookerClientId     *string `gorm:"column:looker_client_id;type:varchar(255);" binding:"required"`
	LookerClientSecret *string `gorm:"column:looker_client_secret;type:varchar(255);" binding:"required"`
	LookerDomain       *string `gorm:"column:looker_domain;type:varchar(255);" binding:"required"`

	MailchimpApiKey *string `gorm:"column:mailchimp_api_key;type:varchar(255);" binding:"required"`

	MailjetEmailApiKey    *string `gorm:"column:mailjet_email_api_key;type:varchar(255);" binding:"required"`
	MailjetEmailApiSecret *string `gorm:"column:mailjet_email_api_secret;type:varchar(255);" binding:"required"`

	MarketoClientId     *string `gorm:"column:marketo_client_id;type:varchar(255);" binding:"required"`
	MarketoClientSecret *string `gorm:"column:marketo_client_secret;type:varchar(255);" binding:"required"`
	MarketoDomainUrl    *string `gorm:"column:marketo_domain_url;type:varchar(255);" binding:"required"`

	MicrosoftTeamsTenantId     *string `gorm:"column:microsoft_teams_tenant_id;type:varchar(255);" binding:"required"`
	MicrosoftTeamsClientId     *string `gorm:"column:microsoft_teams_client_id;type:varchar(255);" binding:"required"`
	MicrosoftTeamsClientSecret *string `gorm:"column:microsoft_teams_client_secret;type:varchar(255);" binding:"required"`

	MondayApiToken *string `gorm:"column:monday_api_token;type:varchar(255);" binding:"required"`

	NotionInternalAccessToken *string `gorm:"column:notion_internal_access_token;type:varchar(255);" binding:"required"`
	NotionPublicAccessToken   *string `gorm:"column:notion_public_access_token;type:varchar(255);" binding:"required"`
	NotionPublicClientId      *string `gorm:"column:notion_public_client_id;type:varchar(255);" binding:"required"`
	NotionPublicClientSecret  *string `gorm:"column:notion_public_client_secret;type:varchar(255);" binding:"required"`

	OracleNetsuiteAccountId      *string `gorm:"column:oracle_netsuite_account_id;type:varchar(255);" binding:"required"`
	OracleNetsuiteConsumerKey    *string `gorm:"column:oracle_netsuite_consumer_key;type:varchar(255);" binding:"required"`
	OracleNetsuiteConsumerSecret *string `gorm:"column:oracle_netsuite_consumer_secret;type:varchar(255);" binding:"required"`
	OracleNetsuiteTokenId        *string `gorm:"column:oracle_netsuite_token_id;type:varchar(255);" binding:"required"`
	OracleNetsuiteTokenSecret    *string `gorm:"column:oracle_netsuite_token_secret;type:varchar(255);" binding:"required"`

	OrbApiKey *string `gorm:"column:orb_api_key;type:varchar(255);" binding:"required"`

	OrbitApiKey *string `gorm:"column:orbit_api_key;type:varchar(255);" binding:"required"`

	PagerDutyApikey *string `gorm:"column:pager_duty_apikey;type:varchar(255);" binding:"required"`

	PaypalTransactionClientId *string `gorm:"column:paypal_transaction_client_id;type:varchar(255);" binding:"required"`
	PaypalTransactionSecret   *string `gorm:"column:paypal_transaction_secret;type:varchar(255);" binding:"required"`

	PaystackSecretKey      *string `gorm:"column:paystack_secret_key;type:varchar(255);" binding:"required"`
	PaystackLookbackWindow *string `gorm:"column:paystack_lookback_window;type:varchar(255);" binding:"required"`

	PendoApiToken *string `gorm:"column:pendo_api_token;type:varchar(255);" binding:"required"`

	PipedriveApiToken *string `gorm:"column:pipedrive_api_token;type:varchar(255);" binding:"required"`

	PlaidAccessToken *string `gorm:"column:plaid_access_token;type:varchar(255);" binding:"required"`

	PlausibleApiKey *string `gorm:"column:plausible_api_key;type:varchar(255);" binding:"required"`
	PlausibleSiteId *string `gorm:"column:plausible_site_id;type:varchar(255);" binding:"required"`

	PostHogApiKey  *string `gorm:"column:post_hog_api_key;type:varchar(255);" binding:"required"`
	PostHogBaseUrl *string `gorm:"column:post_hog_base_url;type:varchar(255);" binding:"required"`

	QualarooApiKey   *string `gorm:"column:qualaroo_api_key;type:varchar(255);" binding:"required"`
	QualarooApiToken *string `gorm:"column:qualaroo_api_token;type:varchar(255);" binding:"required"`

	QuickBooksClientId     *string `gorm:"column:quick_books_client_id;type:varchar(255);" binding:"required"`
	QuickBooksClientSecret *string `gorm:"column:quick_books_client_secret;type:varchar(255);" binding:"required"`
	QuickBooksRealmId      *string `gorm:"column:quick_books_realm_id;type:varchar(255);" binding:"required"`
	QuickBooksRefreshToken *string `gorm:"column:quick_books_refresh_token;type:varchar(255);" binding:"required"`

	RechargeApiToken *string `gorm:"column:recharge_api_token;type:varchar(255);" binding:"required"`

	RecruiteeCompanyId *string `gorm:"column:recruitee_company_id;type:varchar(255);" binding:"required"`
	RecruiteeApiKey    *string `gorm:"column:recruitee_api_key;type:varchar(255);" binding:"required"`

	RecurlyApiKey *string `gorm:"column:recurly_api_key;type:varchar(255);" binding:"required"`

	RetentlyApiToken *string `gorm:"column:retently_api_token;type:varchar(255);" binding:"required"`

	SalesforceClientId     *string `gorm:"column:salesforce_client_id;type:varchar(255);" binding:"required"`
	SalesforceClientSecret *string `gorm:"column:salesforce_client_secret;type:varchar(255);" binding:"required"`
	SalesforceRefreshToken *string `gorm:"column:salesforce_refresh_token;type:varchar(255);" binding:"required"`

	SalesloftApiKey *string `gorm:"column:salesloft_api_key;type:varchar(255);" binding:"required"`

	SendgridApiKey *string `gorm:"column:sendgrid_api_key;type:varchar(255);" binding:"required"`

	SentryProject             *string `gorm:"column:sentry_project;type:varchar(255);" binding:"required"`
	SentryHost                *string `gorm:"column:sentry_host;type:varchar(255);" binding:"required"`
	SentryAuthenticationToken *string `gorm:"column:sentry_authentication_token;type:varchar(255);" binding:"required"`
	SentryOrganization        *string `gorm:"column:sentry_organization;type:varchar(255);" binding:"required"`

	SlackApiToken       *string `gorm:"column:slack_api_token;type:varchar(255);" binding:"required"`
	SlackChannelFilter  *string `gorm:"column:slack_channel_filter;type:varchar(255);" binding:"required"`
	SlackLookbackWindow *string `gorm:"column:slack_lookback_window;type:varchar(255);" binding:"required"`

	StripeAccountId *string `gorm:"column:stripe_account_id;type:varchar(255);" binding:"required"`
	StripeSecretKey *string `gorm:"column:stripe_secret_key;type:varchar(255);" binding:"required"`

	SurveySparrowAccessToken *string `gorm:"column:survey_sparrow_access_token;type:varchar(255);" binding:"required"`

	SurveyMonkeyAccessToken *string `gorm:"column:survey_monkey_access_token;type:varchar(255);" binding:"required"`

	TalkdeskApiKey *string `gorm:"column:talkdesk_api_key;type:varchar(255);" binding:"required"`

	TikTokAccessToken *string `gorm:"column:tik_tok_access_token;type:varchar(255);" binding:"required"`

	TodoistApiToken *string `gorm:"column:todoist_api_token;type:varchar(255);" binding:"required"`

	TypeformApiToken *string `gorm:"column:typeform_api_token;type:varchar(255);" binding:"required"`

	VittallyApiKey *string `gorm:"column:vittally_api_key;type:varchar(255);" binding:"required"`

	WrikeAccessToken *string `gorm:"column:wrike_access_token;type:varchar(255);" binding:"required"`
	WrikeHostUrl     *string `gorm:"column:wrike_host_url;type:varchar(255);" binding:"required"`

	XeroClientId     *string `gorm:"column:xero_client_id;type:varchar(255);" binding:"required"`
	XeroClientSecret *string `gorm:"column:xero_client_secret;type:varchar(255);" binding:"required"`
	XeroTenantId     *string `gorm:"column:xero_tenant_id;type:varchar(255);" binding:"required"`
	XeroScopes       *string `gorm:"column:xero_scopes;type:varchar(255);" binding:"required"`

	ZendeskAPIKey     *string `gorm:"column:zendesk_api_key;type:varchar(255);" binding:"required"`
	ZendeskSubdomain  *string `gorm:"column:zendesk_subdomain;type:varchar(255);" binding:"required"`
	ZendeskAdminEmail *string `gorm:"column:zendesk_admin_email;type:varchar(255);" binding:"required"`

	ZendeskChatSubdomain *string `gorm:"column:zendesk_chat_subdomain;type:varchar(255);" binding:"required"`
	ZendeskChatAccessKey *string `gorm:"column:zendesk_chat_access_key;type:varchar(255);" binding:"required"`

	ZendeskTalkSubdomain *string `gorm:"column:zendesk_talk_subdomain;type:varchar(255);" binding:"required"`
	ZendeskTalkAccessKey *string `gorm:"column:zendesk_talk_access_key;type:varchar(255);" binding:"required"`

	ZendeskSellApiToken *string `gorm:"column:zendesk_sell_api_token;type:varchar(255);" binding:"required"`

	ZendeskSunshineSubdomain *string `gorm:"column:zendesk_sunshine_subdomain;type:varchar(255);" binding:"required"`
	ZendeskSunshineApiToken  *string `gorm:"column:zendesk_sunshine_api_token;type:varchar(255);" binding:"required"`
	ZendeskSunshineEmail     *string `gorm:"column:zendesk_sunshine_email;type:varchar(255);" binding:"required"`

	ZenefitsToken *string `gorm:"column:zenefits_token;type:varchar(255);" binding:"required"`

	MixpanelUsername        *string `gorm:"column:mixpanel_username;type:varchar(255);" binding:"required"`
	MixpanelSecret          *string `gorm:"column:mixpanel_secret;type:varchar(255);" binding:"required"`
	MixpanelProjectId       *string `gorm:"column:mixpanel_project_id;type:varchar(255);" binding:"required"`
	MixpanelProjectSecret   *string `gorm:"column:mixpanel_project_secret;type:varchar(255);" binding:"required"`
	MixpanelProjectTimezone *string `gorm:"column:mixpanel_project_timezone;type:varchar(255);" binding:"required"`
	MixpanelRegion          *string `gorm:"column:mixpanel_region;type:varchar(255);" binding:"required"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}
