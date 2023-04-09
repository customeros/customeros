package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

// TODO the state should come from the actual running service
func MapTenantSettingsEntityToDTO(tenantSettings *entity.TenantSettings) *map[string]interface{} {
	responseMap := make(map[string]interface{})

	if tenantSettings == nil {
		return &responseMap
	}

	if tenantSettings.HubspotPrivateAppKey != nil {
		responseMap["hubspot"] = make(map[string]interface{})
		responseMap["hubspot"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskAPIKey != nil && tenantSettings.ZendeskSubdomain != nil && tenantSettings.ZendeskAdminEmail != nil {
		responseMap["zendesk"] = make(map[string]interface{})
		responseMap["zendesk"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SmartSheetId != nil && tenantSettings.SmartSheetAccessToken != nil {
		responseMap["smartsheet"] = make(map[string]interface{})
		responseMap["smartsheet"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.JiraAPIToken != nil && tenantSettings.JiraDomain != nil && tenantSettings.JiraEmail != nil {
		responseMap["jira"] = make(map[string]interface{})
		responseMap["jira"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TrelloAPIToken != nil && tenantSettings.TrelloAPIKey != nil {
		responseMap["trello"] = make(map[string]interface{})
		responseMap["trello"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AhaAPIUrl != nil && tenantSettings.AhaAPIKey != nil {
		responseMap["aha"] = make(map[string]interface{})
		responseMap["aha"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AirtablePersonalAccessToken != nil {
		responseMap["airtable"] = make(map[string]interface{})
		responseMap["airtable"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AmplitudeSecretKey != nil && tenantSettings.AmplitudeAPIKey != nil {
		responseMap["amplitude"] = make(map[string]interface{})
		responseMap["amplitude"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AsanaAccessToken != nil {
		responseMap["asana"] = make(map[string]interface{})
		responseMap["asana"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BatonAPIKey != nil {
		responseMap["baton"] = make(map[string]interface{})
		responseMap["baton"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BabelforceRegionEnvironment != nil && tenantSettings.BabelforceAccessKeyId != nil && tenantSettings.BabelforceAccessToken != nil {
		responseMap["babelforce"] = make(map[string]interface{})
		responseMap["babelforce"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BigQueryServiceAccountKey != nil {
		responseMap["bigquery"] = make(map[string]interface{})
		responseMap["bigquery"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BraintreeEnvironment != nil && tenantSettings.BraintreeMerchantId != nil && tenantSettings.BraintreePublicKey != nil && tenantSettings.BraintreePrivateKey != nil {
		responseMap["braintree"] = make(map[string]interface{})
		responseMap["braintree"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CallRailAccount != nil && tenantSettings.CallRailApiToken != nil {
		responseMap["callrail"] = make(map[string]interface{})
		responseMap["callrail"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargebeeProductCatalog != nil && tenantSettings.ChargebeeApiKey != nil {
		responseMap["chargebee"] = make(map[string]interface{})
		responseMap["chargebee"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargifyApiKey != nil && tenantSettings.ChargifyDomain != nil {
		responseMap["chargify"] = make(map[string]interface{})
		responseMap["chargify"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ClickUpApiKey != nil {
		responseMap["clickup"] = make(map[string]interface{})
		responseMap["clickup"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CloseComApiKey != nil {
		responseMap["closecom"] = make(map[string]interface{})
		responseMap["closecom"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CodaAuthToken != nil && tenantSettings.CodaDocumentId != nil {
		responseMap["coda"] = make(map[string]interface{})
		responseMap["coda"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ConfluenceApiToken != nil && tenantSettings.ConfluenceDomain != nil && tenantSettings.ConfluenceLoginEmail != nil {
		responseMap["confluence"] = make(map[string]interface{})
		responseMap["confluence"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CourierApiKey != nil {
		responseMap["courier"] = make(map[string]interface{})
		responseMap["courier"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CustomerIoApiKey != nil {
		responseMap["customerio"] = make(map[string]interface{})
		responseMap["customerio"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DatadogApiKey != nil && tenantSettings.DatadogApplicationKey != nil {
		responseMap["datadog"] = make(map[string]interface{})
		responseMap["datadog"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DelightedApiKey != nil {
		responseMap["delighted"] = make(map[string]interface{})
		responseMap["delighted"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DixaApiToken != nil {
		responseMap["dixa"] = make(map[string]interface{})
		responseMap["dixa"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DriftApiToken != nil {
		responseMap["drift"] = make(map[string]interface{})
		responseMap["drift"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.EmailOctopusApiKey != nil {
		responseMap["emailoctopus"] = make(map[string]interface{})
		responseMap["emailoctopus"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FacebookMarketingAccessToken != nil {
		responseMap["facebookMarketing"] = make(map[string]interface{})
		responseMap["facebookMarketing"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FastbillApiKey != nil && tenantSettings.FastbillProjectId != nil {
		responseMap["fastbill"] = make(map[string]interface{})
		responseMap["fastbill"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FlexportApiKey != nil {
		responseMap["flexport"] = make(map[string]interface{})
		responseMap["flexport"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshcallerApiKey != nil {
		responseMap["freshcaller"] = make(map[string]interface{})
		responseMap["freshcaller"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshdeskApiKey != nil && tenantSettings.FreshdeskDomain != nil {
		responseMap["freshdesk"] = make(map[string]interface{})
		responseMap["freshdesk"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshsalesApiKey != nil && tenantSettings.FreshsalesDomain != nil {
		responseMap["freshsales"] = make(map[string]interface{})
		responseMap["freshsales"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshserviceApiKey != nil && tenantSettings.FreshserviceDomain != nil {
		responseMap["freshservice"] = make(map[string]interface{})
		responseMap["freshservice"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GenesysRegion != nil && tenantSettings.GenesysClientId != nil && tenantSettings.GenesysClientSecret != nil {
		responseMap["genesys"] = make(map[string]interface{})
		responseMap["genesys"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitHubAccessToken != nil {
		responseMap["github"] = make(map[string]interface{})
		responseMap["github"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitLabAccessToken != nil {
		responseMap["gitlab"] = make(map[string]interface{})
		responseMap["gitlab"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GoCardlessAccessToken != nil && tenantSettings.GoCardlessEnvironment != nil && tenantSettings.GoCardlessVersion != nil {
		responseMap["gocardless"] = make(map[string]interface{})
		responseMap["gocardless"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GongApiKey != nil {
		responseMap["gong"] = make(map[string]interface{})
		responseMap["gong"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.HarvestAccountId != nil && tenantSettings.HarvestAccessToken != nil {
		responseMap["harvest"] = make(map[string]interface{})
		responseMap["harvest"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InsightlyApiToken != nil {
		responseMap["insightly"] = make(map[string]interface{})
		responseMap["insightly"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstagramAccessToken != nil {
		responseMap["instagram"] = make(map[string]interface{})
		responseMap["instagram"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstatusApiKey != nil {
		responseMap["instatus"] = make(map[string]interface{})
		responseMap["instatus"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.IntercomAccessToken != nil {
		responseMap["intercom"] = make(map[string]interface{})
		responseMap["intercom"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KlaviyoApiKey != nil {
		responseMap["klaviyo"] = make(map[string]interface{})
		responseMap["klaviyo"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KustomerApiToken != nil {
		responseMap["kustomer"] = make(map[string]interface{})
		responseMap["kustomer"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.LookerClientId != nil && tenantSettings.LookerClientSecret != nil && tenantSettings.LookerDomain != nil {
		responseMap["looker"] = make(map[string]interface{})
		responseMap["looker"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailchimpApiKey != nil {
		responseMap["mailchimp"] = make(map[string]interface{})
		responseMap["mailchimp"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailjetEmailApiKey != nil && tenantSettings.MailjetEmailApiSecret != nil {
		responseMap["mailjetemail"] = make(map[string]interface{})
		responseMap["mailjetemail"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MarketoClientId != nil && tenantSettings.MarketoClientSecret != nil && tenantSettings.MarketoDomainUrl != nil {
		responseMap["marketo"] = make(map[string]interface{})
		responseMap["marketo"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MicrosoftTeamsTenantId != nil && tenantSettings.MicrosoftTeamsClientId != nil && tenantSettings.MicrosoftTeamsClientSecret != nil {
		responseMap["microsoftteams"] = make(map[string]interface{})
		responseMap["microsoftteams"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MondayApiToken != nil {
		responseMap["monday"] = make(map[string]interface{})
		responseMap["monday"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.NotionInternalAccessToken != nil || (tenantSettings.NotionPublicClientId != nil && tenantSettings.NotionPublicClientSecret != nil && tenantSettings.NotionPublicAccessToken != nil) {
		responseMap["notion"] = make(map[string]interface{})
		responseMap["notion"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbApiKey != nil {
		responseMap["orb"] = make(map[string]interface{})
		responseMap["orb"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbitApiKey != nil {
		responseMap["orbit"] = make(map[string]interface{})
		responseMap["orbit"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PagerDutyApikey != nil {
		responseMap["pagerduty"] = make(map[string]interface{})
		responseMap["pagerduty"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaypalTransactionClientId != nil && tenantSettings.PaypalTransactionSecret != nil {
		responseMap["paypaltransaction"] = make(map[string]interface{})
		responseMap["paypaltransaction"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaystackSecretKey != nil {
		responseMap["paystack"] = make(map[string]interface{})
		responseMap["paystack"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PendoApiToken != nil {
		responseMap["pendo"] = make(map[string]interface{})
		responseMap["pendo"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PipedriveApiToken != nil {
		responseMap["pipedrive"] = make(map[string]interface{})
		responseMap["pipedrive"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlaidAccessToken != nil {
		responseMap["plaid"] = make(map[string]interface{})
		responseMap["plaid"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlausibleSiteId != nil && tenantSettings.PlausibleApiKey != nil {
		responseMap["plausible"] = make(map[string]interface{})
		responseMap["plausible"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PostHogApiKey != nil {
		responseMap["posthog"] = make(map[string]interface{})
		responseMap["posthog"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QualarooApiKey != nil && tenantSettings.QualarooApiToken != nil {
		responseMap["qualaroo"] = make(map[string]interface{})
		responseMap["qualaroo"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QuickBooksClientId != nil && tenantSettings.QuickBooksClientSecret != nil && tenantSettings.QuickBooksRealmId != nil && tenantSettings.QuickBooksRefreshToken != nil {
		responseMap["quickbooks"] = make(map[string]interface{})
		responseMap["quickbooks"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RechargeApiToken != nil {
		responseMap["recharge"] = make(map[string]interface{})
		responseMap["recharge"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecruiteeCompanyId != nil && tenantSettings.RecruiteeApiKey != nil {
		responseMap["recruitee"] = make(map[string]interface{})
		responseMap["recruitee"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecurlyApiKey != nil {
		responseMap["recurly"] = make(map[string]interface{})
		responseMap["recurly"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RetentlyApiToken != nil {
		responseMap["retently"] = make(map[string]interface{})
		responseMap["retently"].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SalesloftApiKey != nil {
		responseMap["salesloft"] = make(map[string]interface{})
		responseMap["salesloft"].(map[string]interface{})["state"] = "ACTIVE"
	}

	return &responseMap
}
