package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func hasMixpanelKeys(tenantSettings *entity.TenantSettings) bool {
	return tenantSettings.MixpanelUsername != nil || tenantSettings.MixpanelSecret != nil || tenantSettings.MixpanelProjectId != nil || tenantSettings.MixpanelProjectSecret != nil || tenantSettings.MixpanelProjectTimezone != nil || tenantSettings.MixpanelRegion != nil
}

// TODO the state should come from the actual running service
func MapTenantSettingsEntityToDTO(tenantSettings *entity.TenantSettings, activeServices map[string]bool) *map[string]interface{} {
	responseMap := make(map[string]interface{})

	for service, isActive := range activeServices {
		if isActive {
			responseMap[service] = make(map[string]interface{})
			responseMap[service].(map[string]interface{})["state"] = "ACTIVE"
		}
	}

	if tenantSettings == nil {
		return &responseMap
	}

	if tenantSettings.HubspotPrivateAppKey != nil {
		responseMap[service.SERVICE_HUBSPOT] = make(map[string]interface{})
		responseMap[service.SERVICE_HUBSPOT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SmartSheetId != nil && tenantSettings.SmartSheetAccessToken != nil {
		responseMap[service.SERVICE_SMARTSHEET] = make(map[string]interface{})
		responseMap[service.SERVICE_SMARTSHEET].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.JiraAPIToken != nil && tenantSettings.JiraDomain != nil && tenantSettings.JiraEmail != nil {
		responseMap[service.SERVICE_JIRA] = make(map[string]interface{})
		responseMap[service.SERVICE_JIRA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TrelloAPIToken != nil && tenantSettings.TrelloAPIKey != nil {
		responseMap[service.SERVICE_TRELLO] = make(map[string]interface{})
		responseMap[service.SERVICE_TRELLO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AhaAPIUrl != nil && tenantSettings.AhaAPIKey != nil {
		responseMap[service.SERVICE_AHA] = make(map[string]interface{})
		responseMap[service.SERVICE_AHA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AirtablePersonalAccessToken != nil {
		responseMap[service.SERVICE_AIRTABLE] = make(map[string]interface{})
		responseMap[service.SERVICE_AIRTABLE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AmplitudeSecretKey != nil && tenantSettings.AmplitudeAPIKey != nil {
		responseMap[service.SERVICE_AMPLITUDE] = make(map[string]interface{})
		responseMap[service.SERVICE_AMPLITUDE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AsanaAccessToken != nil {
		responseMap[service.SERVICE_ASANA] = make(map[string]interface{})
		responseMap[service.SERVICE_ASANA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BatonAPIKey != nil {
		responseMap[service.SERVICE_BATON] = make(map[string]interface{})
		responseMap[service.SERVICE_BATON].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BabelforceRegionEnvironment != nil && tenantSettings.BabelforceAccessKeyId != nil && tenantSettings.BabelforceAccessToken != nil {
		responseMap[service.SERVICE_BABELFORCE] = make(map[string]interface{})
		responseMap[service.SERVICE_BABELFORCE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BigQueryServiceAccountKey != nil {
		responseMap[service.SERVICE_BIGQUERY] = make(map[string]interface{})
		responseMap[service.SERVICE_BIGQUERY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BraintreeEnvironment != nil && tenantSettings.BraintreeMerchantId != nil && tenantSettings.BraintreePublicKey != nil && tenantSettings.BraintreePrivateKey != nil {
		responseMap[service.SERVICE_BRAINTREE] = make(map[string]interface{})
		responseMap[service.SERVICE_BRAINTREE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CallRailAccount != nil && tenantSettings.CallRailApiToken != nil {
		responseMap[service.SERVICE_CALLRAIL] = make(map[string]interface{})
		responseMap[service.SERVICE_CALLRAIL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargebeeProductCatalog != nil && tenantSettings.ChargebeeApiKey != nil {
		responseMap[service.SERVICE_CHARGEBEE] = make(map[string]interface{})
		responseMap[service.SERVICE_CHARGEBEE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargifyApiKey != nil && tenantSettings.ChargifyDomain != nil {
		responseMap[service.SERVICE_CHARGIFY] = make(map[string]interface{})
		responseMap[service.SERVICE_CHARGIFY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ClickUpApiKey != nil {
		responseMap[service.SERVICE_CLICKUP] = make(map[string]interface{})
		responseMap[service.SERVICE_CLICKUP].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CloseComApiKey != nil {
		responseMap[service.SERVICE_CLOSECOM] = make(map[string]interface{})
		responseMap[service.SERVICE_CLOSECOM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CodaAuthToken != nil && tenantSettings.CodaDocumentId != nil {
		responseMap[service.SERVICE_CODA] = make(map[string]interface{})
		responseMap[service.SERVICE_CODA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ConfluenceApiToken != nil && tenantSettings.ConfluenceDomain != nil && tenantSettings.ConfluenceLoginEmail != nil {
		responseMap[service.SERVICE_CONFLUENCE] = make(map[string]interface{})
		responseMap[service.SERVICE_CONFLUENCE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CourierApiKey != nil {
		responseMap[service.SERVICE_COURIER] = make(map[string]interface{})
		responseMap[service.SERVICE_COURIER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CustomerIoApiKey != nil {
		responseMap[service.SERVICE_CUSTOMERIO] = make(map[string]interface{})
		responseMap[service.SERVICE_CUSTOMERIO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DatadogApiKey != nil && tenantSettings.DatadogApplicationKey != nil {
		responseMap[service.SERVICE_DATADOG] = make(map[string]interface{})
		responseMap[service.SERVICE_DATADOG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DelightedApiKey != nil {
		responseMap[service.SERVICE_DELIGHTED] = make(map[string]interface{})
		responseMap[service.SERVICE_DELIGHTED].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DixaApiToken != nil {
		responseMap[service.SERVICE_DIXA] = make(map[string]interface{})
		responseMap[service.SERVICE_DIXA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DriftApiToken != nil {
		responseMap[service.SERVICE_DRIFT] = make(map[string]interface{})
		responseMap[service.SERVICE_DRIFT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.EmailOctopusApiKey != nil {
		responseMap[service.SERVICE_EMAILOCTOPUS] = make(map[string]interface{})
		responseMap[service.SERVICE_EMAILOCTOPUS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FacebookMarketingAccessToken != nil {
		responseMap[service.SERVICE_FACEBOOK_MARKETING] = make(map[string]interface{})
		responseMap[service.SERVICE_FACEBOOK_MARKETING].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FastbillApiKey != nil && tenantSettings.FastbillProjectId != nil {
		responseMap[service.SERVICE_FASTBILL] = make(map[string]interface{})
		responseMap[service.SERVICE_FASTBILL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FlexportApiKey != nil {
		responseMap[service.SERVICE_FLEXPORT] = make(map[string]interface{})
		responseMap[service.SERVICE_FLEXPORT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshcallerApiKey != nil {
		responseMap[service.SERVICE_FRESHCALLER] = make(map[string]interface{})
		responseMap[service.SERVICE_FRESHCALLER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshdeskApiKey != nil && tenantSettings.FreshdeskDomain != nil {
		responseMap[service.SERVICE_FRESHDESK] = make(map[string]interface{})
		responseMap[service.SERVICE_FRESHDESK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshsalesApiKey != nil && tenantSettings.FreshsalesDomain != nil {
		responseMap[service.SERVICE_FRESHSALES] = make(map[string]interface{})
		responseMap[service.SERVICE_FRESHSALES].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshserviceApiKey != nil && tenantSettings.FreshserviceDomain != nil {
		responseMap[service.SERVICE_FRESHSERVICE] = make(map[string]interface{})
		responseMap[service.SERVICE_FRESHSERVICE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GenesysRegion != nil && tenantSettings.GenesysClientId != nil && tenantSettings.GenesysClientSecret != nil {
		responseMap[service.SERVICE_GENESYS] = make(map[string]interface{})
		responseMap[service.SERVICE_GENESYS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitHubAccessToken != nil {
		responseMap[service.SERVICE_GITHUB] = make(map[string]interface{})
		responseMap[service.SERVICE_GITHUB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitLabAccessToken != nil {
		responseMap[service.SERVICE_GITLAB] = make(map[string]interface{})
		responseMap[service.SERVICE_GITLAB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GoCardlessAccessToken != nil && tenantSettings.GoCardlessEnvironment != nil && tenantSettings.GoCardlessVersion != nil {
		responseMap[service.SERVICE_GOCARDLESS] = make(map[string]interface{})
		responseMap[service.SERVICE_GOCARDLESS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GongApiKey != nil {
		responseMap[service.SERVICE_GONG] = make(map[string]interface{})
		responseMap[service.SERVICE_GONG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.HarvestAccountId != nil && tenantSettings.HarvestAccessToken != nil {
		responseMap[service.SERVICE_HARVEST] = make(map[string]interface{})
		responseMap[service.SERVICE_HARVEST].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InsightlyApiToken != nil {
		responseMap[service.SERVICE_INSIGHTLY] = make(map[string]interface{})
		responseMap[service.SERVICE_INSIGHTLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstagramAccessToken != nil {
		responseMap[service.SERVICE_INSTAGRAM] = make(map[string]interface{})
		responseMap[service.SERVICE_INSTAGRAM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstatusApiKey != nil {
		responseMap[service.SERVICE_INSTATUS] = make(map[string]interface{})
		responseMap[service.SERVICE_INSTATUS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.IntercomAccessToken != nil {
		responseMap[service.SERVICE_INTERCOM] = make(map[string]interface{})
		responseMap[service.SERVICE_INTERCOM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KlaviyoApiKey != nil {
		responseMap[service.SERVICE_KLAVIYO] = make(map[string]interface{})
		responseMap[service.SERVICE_KLAVIYO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KustomerApiToken != nil {
		responseMap[service.SERVICE_KUSTOMER] = make(map[string]interface{})
		responseMap[service.SERVICE_KUSTOMER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.LookerClientId != nil && tenantSettings.LookerClientSecret != nil && tenantSettings.LookerDomain != nil {
		responseMap[service.SERVICE_LOOKER] = make(map[string]interface{})
		responseMap[service.SERVICE_LOOKER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailchimpApiKey != nil {
		responseMap[service.SERVICE_MAILCHIMP] = make(map[string]interface{})
		responseMap[service.SERVICE_MAILCHIMP].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailjetEmailApiKey != nil && tenantSettings.MailjetEmailApiSecret != nil {
		responseMap[service.SERVICE_MAILJETEMAIL] = make(map[string]interface{})
		responseMap[service.SERVICE_MAILJETEMAIL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MarketoClientId != nil && tenantSettings.MarketoClientSecret != nil && tenantSettings.MarketoDomainUrl != nil {
		responseMap[service.SERVICE_MARKETO] = make(map[string]interface{})
		responseMap[service.SERVICE_MARKETO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MicrosoftTeamsTenantId != nil && tenantSettings.MicrosoftTeamsClientId != nil && tenantSettings.MicrosoftTeamsClientSecret != nil {
		responseMap[service.SERVICE_MICROSOFT_TEAMS] = make(map[string]interface{})
		responseMap[service.SERVICE_MICROSOFT_TEAMS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MondayApiToken != nil {
		responseMap[service.SERVICE_MONDAY] = make(map[string]interface{})
		responseMap[service.SERVICE_MONDAY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.NotionInternalAccessToken != nil || (tenantSettings.NotionPublicClientId != nil && tenantSettings.NotionPublicClientSecret != nil && tenantSettings.NotionPublicAccessToken != nil) {
		responseMap[service.SERVICE_NOTION] = make(map[string]interface{})
		responseMap[service.SERVICE_NOTION].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OracleNetsuiteAccountId != nil && tenantSettings.OracleNetsuiteConsumerKey != nil && tenantSettings.OracleNetsuiteConsumerSecret != nil && tenantSettings.OracleNetsuiteTokenId != nil && tenantSettings.OracleNetsuiteTokenSecret != nil {
		responseMap[service.SERVICE_ORACLE_NETSUITE] = make(map[string]interface{})
		responseMap[service.SERVICE_ORACLE_NETSUITE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbApiKey != nil {
		responseMap[service.SERVICE_ORB] = make(map[string]interface{})
		responseMap[service.SERVICE_ORB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbitApiKey != nil {
		responseMap[service.SERVICE_ORBIT] = make(map[string]interface{})
		responseMap[service.SERVICE_ORBIT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PagerDutyApikey != nil {
		responseMap[service.SERVICE_PAGERDUTY] = make(map[string]interface{})
		responseMap[service.SERVICE_PAGERDUTY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaypalTransactionClientId != nil && tenantSettings.PaypalTransactionSecret != nil {
		responseMap[service.SERVICE_PAYPAL_TRANSACTION] = make(map[string]interface{})
		responseMap[service.SERVICE_PAYPAL_TRANSACTION].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaystackSecretKey != nil {
		responseMap[service.SERVICE_PAYSTACK] = make(map[string]interface{})
		responseMap[service.SERVICE_PAYSTACK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PendoApiToken != nil {
		responseMap[service.SERVICE_PENDO] = make(map[string]interface{})
		responseMap[service.SERVICE_PENDO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PipedriveApiToken != nil {
		responseMap[service.SERVICE_PIPEDRIVE] = make(map[string]interface{})
		responseMap[service.SERVICE_PIPEDRIVE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlaidAccessToken != nil {
		responseMap[service.SERVICE_PLAID] = make(map[string]interface{})
		responseMap[service.SERVICE_PLAID].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlausibleSiteId != nil && tenantSettings.PlausibleApiKey != nil {
		responseMap[service.SERVICE_PLAUSIBLE] = make(map[string]interface{})
		responseMap[service.SERVICE_PLAUSIBLE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PostHogApiKey != nil {
		responseMap[service.SERVICE_POSTHOG] = make(map[string]interface{})
		responseMap[service.SERVICE_POSTHOG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QualarooApiKey != nil && tenantSettings.QualarooApiToken != nil {
		responseMap[service.SERVICE_QUALAROO] = make(map[string]interface{})
		responseMap[service.SERVICE_QUALAROO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QuickBooksClientId != nil && tenantSettings.QuickBooksClientSecret != nil && tenantSettings.QuickBooksRealmId != nil && tenantSettings.QuickBooksRefreshToken != nil {
		responseMap[service.SERVICE_QUICKBOOKS] = make(map[string]interface{})
		responseMap[service.SERVICE_QUICKBOOKS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RechargeApiToken != nil {
		responseMap[service.SERVICE_RECHARGE] = make(map[string]interface{})
		responseMap[service.SERVICE_RECHARGE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecruiteeCompanyId != nil && tenantSettings.RecruiteeApiKey != nil {
		responseMap[service.SERVICE_RECRUITEE] = make(map[string]interface{})
		responseMap[service.SERVICE_RECRUITEE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecurlyApiKey != nil {
		responseMap[service.SERVICE_RECURLY] = make(map[string]interface{})
		responseMap[service.SERVICE_RECURLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RetentlyApiToken != nil {
		responseMap[service.SERVICE_RETENTLY] = make(map[string]interface{})
		responseMap[service.SERVICE_RETENTLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SalesforceClientId != nil && tenantSettings.SalesforceClientSecret != nil && tenantSettings.SalesforceRefreshToken != nil {
		responseMap[service.SERVICE_SALESFORCE] = make(map[string]interface{})
		responseMap[service.SERVICE_SALESFORCE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SalesloftApiKey != nil {
		responseMap[service.SERVICE_SALESLOFT] = make(map[string]interface{})
		responseMap[service.SERVICE_SALESLOFT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SendgridApiKey != nil {
		responseMap[service.SERVICE_SENDGRID] = make(map[string]interface{})
		responseMap[service.SERVICE_SENDGRID].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SentryProject != nil && tenantSettings.SentryOrganization != nil && tenantSettings.SentryAuthenticationToken != nil {
		responseMap[service.SERVICE_SENTRY] = make(map[string]interface{})
		responseMap[service.SERVICE_SENTRY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SlackApiToken != nil && tenantSettings.SlackChannelFilter != nil {
		responseMap[service.SERVICE_SLACK] = make(map[string]interface{})
		responseMap[service.SERVICE_SLACK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.StripeAccountId != nil && tenantSettings.StripeSecretKey != nil {
		responseMap[service.SERVICE_STRIPE] = make(map[string]interface{})
		responseMap[service.SERVICE_STRIPE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SurveySparrowAccessToken != nil {
		responseMap[service.SERVICE_SURVEYSPARROW] = make(map[string]interface{})
		responseMap[service.SERVICE_SURVEYSPARROW].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SurveyMonkeyAccessToken != nil {
		responseMap[service.SERVICE_SURVEYMONKEY] = make(map[string]interface{})
		responseMap[service.SERVICE_SURVEYMONKEY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TalkdeskApiKey != nil {
		responseMap[service.SERVICE_TALKDESK] = make(map[string]interface{})
		responseMap[service.SERVICE_TALKDESK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TikTokAccessToken != nil {
		responseMap[service.SERVICE_TIKTOK] = make(map[string]interface{})
		responseMap[service.SERVICE_TIKTOK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TodoistApiToken != nil {
		responseMap[service.SERVICE_TODOIST] = make(map[string]interface{})
		responseMap[service.SERVICE_TODOIST].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TypeformApiToken != nil {
		responseMap[service.SERVICE_TYPEFORM] = make(map[string]interface{})
		responseMap[service.SERVICE_TYPEFORM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.VittallyApiKey != nil {
		responseMap[service.SERVICE_VITTALLY] = make(map[string]interface{})
		responseMap[service.SERVICE_VITTALLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.WrikeAccessToken != nil && tenantSettings.WrikeHostUrl != nil {
		responseMap[service.SERVICE_WRIKE] = make(map[string]interface{})
		responseMap[service.SERVICE_WRIKE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.XeroClientId != nil && tenantSettings.XeroClientSecret != nil && tenantSettings.XeroTenantId != nil && tenantSettings.XeroScopes != nil {
		responseMap[service.SERVICE_XERO] = make(map[string]interface{})
		responseMap[service.SERVICE_XERO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskAPIKey != nil && tenantSettings.ZendeskSubdomain != nil && tenantSettings.ZendeskAdminEmail != nil {
		responseMap[service.SERVICE_ZENDESK_SUPPORT] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENDESK_SUPPORT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskChatSubdomain != nil && tenantSettings.ZendeskChatAccessKey != nil {
		responseMap[service.SERVICE_ZENDESK_CHAT] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENDESK_CHAT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskTalkSubdomain != nil && tenantSettings.ZendeskTalkAccessKey != nil {
		responseMap[service.SERVICE_ZENDESK_TALK] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENDESK_TALK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskSellApiToken != nil {
		responseMap[service.SERVICE_ZENDESK_SELL] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENDESK_SELL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskSunshineSubdomain != nil && tenantSettings.ZendeskSunshineApiToken != nil && tenantSettings.ZendeskSunshineEmail != nil {
		responseMap[service.SERVICE_ZENDESK_SUNSHINE] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENDESK_SUNSHINE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZenefitsToken != nil {
		responseMap[service.SERVICE_ZENEFITS] = make(map[string]interface{})
		responseMap[service.SERVICE_ZENEFITS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && hasMixpanelKeys(tenantSettings) {
		responseMap[service.SERVICE_MIXPANEL] = make(map[string]interface{})
		responseMap[service.SERVICE_MIXPANEL].(map[string]interface{})["state"] = "ACTIVE"
	}

	return &responseMap
}
