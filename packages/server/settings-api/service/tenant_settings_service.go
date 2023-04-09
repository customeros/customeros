package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

type TenantSettingsService interface {
	GetForTenant(tenantName string) (*entity.TenantSettings, error)

	SaveIntegrationData(tenantName string, request map[string]interface{}) (*entity.TenantSettings, error)
	ClearIntegrationData(tenantName, identifier string) (*entity.TenantSettings, error)
}

type tenantSettingsService struct {
	repositories *repository.PostgresRepositories
}

func NewTenantSettingsService(repositories *repository.PostgresRepositories) TenantSettingsService {
	return &tenantSettingsService{
		repositories: repositories,
	}
}

func (s *tenantSettingsService) GetForTenant(tenantName string) (*entity.TenantSettings, error) {
	qr := s.repositories.TenantSettingsRepository.FindForTenantName(tenantName)
	if qr.Error != nil {
		return nil, qr.Error
	} else if qr.Result == nil {
		return nil, nil
	} else {
		settings := qr.Result.(entity.TenantSettings)
		return &settings, nil
	}
}

func (s *tenantSettingsService) SaveIntegrationData(tenantName string, request map[string]interface{}) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		tenantSettings = &entity.TenantSettings{
			TenantName: tenantName,
		}

		if qr := s.repositories.TenantSettingsRepository.Save(tenantSettings); qr.Error != nil {
			return nil, qr.Error
		}
	}

	// Update tenant settings with new integration data
	for integrationId, value := range request {
		data, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data for integration %s", integrationId)
		}

		switch integrationId {
		case "hubspot":
			privateAppKey, ok := data["privateAppKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing private app key for Hubspot integration")
			}
			tenantSettings.HubspotPrivateAppKey = &privateAppKey

		case "zendesk":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Zendesk integration")
			}
			subdomain, ok := data["subdomain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing subdomain for Zendesk integration")
			}
			adminEmail, ok := data["adminEmail"].(string)
			if !ok {
				return nil, fmt.Errorf("missing admin email for Zendesk integration")
			}
			tenantSettings.ZendeskAPIKey = &apiKey
			tenantSettings.ZendeskSubdomain = &subdomain
			tenantSettings.ZendeskAdminEmail = &adminEmail

		case "smartsheet":
			id, ok := data["id"].(string)
			if !ok {
				return nil, fmt.Errorf("missing Smartsheet ID")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Smartsheet integration")
			}
			tenantSettings.SmartSheetId = &id
			tenantSettings.SmartSheetAccessToken = &accessToken

		case "jira":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Jira integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Jira integration")
			}
			email, ok := data["email"].(string)
			if !ok {
				return nil, fmt.Errorf("missing email for Jira integration")
			}
			tenantSettings.JiraAPIToken = &apiToken
			tenantSettings.JiraDomain = &domain
			tenantSettings.JiraEmail = &email

		case "trello":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Trello integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Trello integration")
			}
			tenantSettings.TrelloAPIToken = &apiToken
			tenantSettings.TrelloAPIKey = &apiKey

		case "aha":
			apiUrl, ok := data["apiUrl"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API Url for Aha integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Aha integration")
			}
			tenantSettings.AhaAPIUrl = &apiUrl
			tenantSettings.AhaAPIKey = &apiKey

		case "airtable":
			personalAccessToken, ok := data["personalAccessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing personal access token for Airtable integration")
			}
			tenantSettings.AirtablePersonalAccessToken = &personalAccessToken

		case "amplitude":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Amplitude integration")
			}
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing secret key for Amplitude integration")
			}
			tenantSettings.AmplitudeSecretKey = &secretKey
			tenantSettings.AmplitudeAPIKey = &apiKey

		case "asana":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Asana integration")
			}

			tenantSettings.AsanaAccessToken = &accessToken

		case "baton":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Baton integration")
			}
			tenantSettings.BatonAPIKey = &apiKey

		case "babelforce":
			regionEnvironment, ok := data["regionEnvironment"].(string)
			if !ok {
				return nil, fmt.Errorf("missing region / environment for Babelforce integration")
			}
			accessKeyId, ok := data["accessKeyId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access key id for Babelforce integration")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Babelforce integration")
			}

			tenantSettings.BabelforceRegionEnvironment = &regionEnvironment
			tenantSettings.BabelforceAccessKeyId = &accessKeyId
			tenantSettings.BabelforceAccessToken = &accessToken

		case "bigquery":
			serviceAccountKey, ok := data["serviceAccountKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing service account key for BigQuery integration")
			}

			tenantSettings.BigQueryServiceAccountKey = &serviceAccountKey

		case "braintree":
			publicKey, ok := data["publicKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing public key for Braintree integration")
			}
			privateKey, ok := data["privateKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing private key for Braintree integration")
			}
			environment, ok := data["environment"].(string)
			if !ok {
				return nil, fmt.Errorf("missing environment for Braintree integration")
			}
			merchantId, ok := data["merchantId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing merchant id for Braintree integration")
			}

			tenantSettings.BraintreePublicKey = &publicKey
			tenantSettings.BraintreePrivateKey = &privateKey
			tenantSettings.BraintreeEnvironment = &environment
			tenantSettings.BraintreeMerchantId = &merchantId

		case "callrail":
			account, ok := data["account"].(string)
			if !ok {
				return nil, fmt.Errorf("missing account for CallRail integration")
			}
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for CallRail integration")
			}

			tenantSettings.CallRailAccount = &account
			tenantSettings.CallRailApiToken = &apiToken

		case "chargebee":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Chargebee integration")
			}
			productCatalog, ok := data["productCatalog"].(string)
			if !ok {
				return nil, fmt.Errorf("missing product catalog for CallRail integration")
			}

			tenantSettings.ChargebeeApiKey = &apiKey
			tenantSettings.ChargebeeProductCatalog = &productCatalog

		case "chargify":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Chargify integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Chargify integration")
			}

			tenantSettings.ChargifyApiKey = &apiKey
			tenantSettings.ChargifyDomain = &domain

		case "clickup":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for ClickUp integration")
			}

			tenantSettings.ClickUpApiKey = &apiKey

		case "closecom":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Close.com integration")
			}

			tenantSettings.CloseComApiKey = &apiKey

		case "coda":
			authToken, ok := data["authToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing auth token for Coda integration")
			}
			documentId, ok := data["documentId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing document id for Coda integration")
			}

			tenantSettings.CodaAuthToken = &authToken
			tenantSettings.CodaDocumentId = &documentId

		case "confluence":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Confluence integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Confluence integration")
			}
			loginEmail, ok := data["loginEmail"].(string)
			if !ok {
				return nil, fmt.Errorf("missing login email for Confluence integration")
			}

			tenantSettings.ConfluenceApiToken = &apiToken
			tenantSettings.ConfluenceDomain = &domain
			tenantSettings.ConfluenceLoginEmail = &loginEmail

		case "courier":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Courier integration")
			}

			tenantSettings.CourierApiKey = &apiKey

		case "customerio":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Customer.io integration")
			}

			tenantSettings.CustomerIoApiKey = &apiKey

		case "datadog":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Customer.io integration")
			}
			applicationKey, ok := data["applicationKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing application key for Customer.io integration")
			}

			tenantSettings.DatadogApiKey = &apiKey
			tenantSettings.DatadogApplicationKey = &applicationKey

		case "delighted":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Delighted integration")
			}

			tenantSettings.DelightedApiKey = &apiKey

		case "dixa":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Dixa integration")
			}

			tenantSettings.DixaApiToken = &apiToken

		case "drift":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Drift integration")
			}

			tenantSettings.DriftApiToken = &apiToken

		case "emailoctopus":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for EmailOctopus integration")
			}

			tenantSettings.EmailOctopusApiKey = &apiKey

		case "facebookMarketing":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Facebook integration")
			}

			tenantSettings.FacebookMarketingAccessToken = &accessToken

		case "fastbill":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Fastbill integration")
			}
			projectId, ok := data["projectId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing project id for Fastbill integration")
			}

			tenantSettings.FastbillApiKey = &apiKey
			tenantSettings.FastbillProjectId = &projectId

		case "flexport":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Flexport integration")
			}

			tenantSettings.FlexportApiKey = &apiKey

		case "freshcaller":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Freshcaller integration")
			}

			tenantSettings.FreshcallerApiKey = &apiKey

		case "freshdesk":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Freshdesk integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Freshdesk integration")
			}

			tenantSettings.FreshdeskApiKey = &apiKey
			tenantSettings.FreshdeskDomain = &domain

		case "freshsales":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Freshsales integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Freshsales integration")
			}

			tenantSettings.FreshsalesApiKey = &apiKey
			tenantSettings.FreshsalesDomain = &domain

		case "freshservice":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Freshservice integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Freshservice integration")
			}

			tenantSettings.FreshserviceApiKey = &apiKey
			tenantSettings.FreshserviceDomain = &domain

		case "genesys":
			region, ok := data["region"].(string)
			if !ok {
				return nil, fmt.Errorf("missing region for Genesys integration")
			}
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for Genesys integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client secret for Genesys integration")
			}

			tenantSettings.GenesysRegion = &region
			tenantSettings.GenesysClientId = &clientId
			tenantSettings.GenesysClientSecret = &clientSecret

		case "github":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for GitHub integration")
			}

			tenantSettings.GitHubAccessToken = &accessToken

		case "gitlab":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for GitLab integration")
			}

			tenantSettings.GitLabAccessToken = &accessToken

		case "gocardless":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for GoCardless integration")
			}
			environment, ok := data["environment"].(string)
			if !ok {
				return nil, fmt.Errorf("missing environment for GoCardless integration")
			}
			version, ok := data["version"].(string)
			if !ok {
				return nil, fmt.Errorf("missing version for GoCardless integration")
			}

			tenantSettings.GoCardlessAccessToken = &accessToken
			tenantSettings.GoCardlessEnvironment = &environment
			tenantSettings.GoCardlessVersion = &version

		case "gong":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Gong integration")
			}

			tenantSettings.GongApiKey = &apiKey

		case "harvest":
			accountId, ok := data["accountId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing account id for Harvest integration")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Harvest integration")
			}

			tenantSettings.HarvestAccountId = &accountId
			tenantSettings.HarvestAccessToken = &accessToken

		case "insightly":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Insightly integration")
			}

			tenantSettings.InsightlyApiToken = &apiToken

		case "instagram":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Harvest integration")
			}

			tenantSettings.InstagramAccessToken = &accessToken

		case "instatus":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Instatus integration")
			}

			tenantSettings.InstatusApiKey = &apiKey

		case "intercom":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Intercom integration")
			}

			tenantSettings.IntercomAccessToken = &accessToken

		case "klaviyo":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Klaviyo integration")
			}

			tenantSettings.KlaviyoApiKey = &apiKey

		case "kustomer":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Kustomer integration")
			}

			tenantSettings.KustomerApiToken = &apiToken

		case "looker":
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for Looker integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client secret for Looker integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain for Looker integration")
			}

			tenantSettings.LookerClientId = &clientId
			tenantSettings.LookerClientSecret = &clientSecret
			tenantSettings.LookerDomain = &domain

		case "mailchimp":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Mailchimp integration")
			}

			tenantSettings.MailchimpApiKey = &apiKey

		case "mailjetemail":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Mailjet Email integration")
			}
			apiSecret, ok := data["apiSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API secret for Mailjet Email integration")
			}

			tenantSettings.MailjetEmailApiKey = &apiKey
			tenantSettings.MailjetEmailApiSecret = &apiSecret

		case "marketo":
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for Marketo integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client secret for Marketo integration")
			}
			domainUrl, ok := data["domainUrl"].(string)
			if !ok {
				return nil, fmt.Errorf("missing domain URL for Marketo integration")
			}

			tenantSettings.MarketoClientId = &clientId
			tenantSettings.MarketoClientSecret = &clientSecret
			tenantSettings.MarketoDomainUrl = &domainUrl

		case "microsoftteams":
			tenantId, ok := data["tenantId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing tenant id for Microsoft Teams integration")
			}
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for Microsoft Teams integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client secret for Microsoft Teams integration")
			}

			tenantSettings.MicrosoftTeamsTenantId = &tenantId
			tenantSettings.MicrosoftTeamsClientId = &clientId
			tenantSettings.MicrosoftTeamsClientSecret = &clientSecret

		case "monday":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Monday integration")
			}

			tenantSettings.MondayApiToken = &apiToken

		case "notion":
			internalAccessToken, _ := data["internalAccessToken"].(string)

			publicClientId, _ := data["publicClientId"].(string)
			publicClientSecret, _ := data["publicClientSecret"].(string)
			publicAccessToken, _ := data["publicAccessToken"].(string)

			if internalAccessToken == "" && publicClientId == "" && publicClientSecret == "" && publicAccessToken == "" {
				return nil, fmt.Errorf("missing Notion integration data")
			}

			if internalAccessToken == "" && (publicClientId == "" || publicClientSecret == "" || publicAccessToken == "") {
				return nil, fmt.Errorf("missing public Notion integration data")
			}

			if internalAccessToken == "" {
				tenantSettings.NotionInternalAccessToken = nil
			} else {
				tenantSettings.NotionInternalAccessToken = &internalAccessToken
			}

			if publicClientId == "" {
				tenantSettings.NotionPublicClientId = nil
			} else {
				tenantSettings.NotionPublicClientId = &publicClientId
			}

			if publicClientSecret == "" {
				tenantSettings.NotionPublicClientSecret = nil
			} else {
				tenantSettings.NotionPublicClientSecret = &publicClientSecret
			}

			if publicAccessToken == "" {
				tenantSettings.NotionPublicAccessToken = nil
			} else {
				tenantSettings.NotionPublicAccessToken = &publicAccessToken
			}

		case "orb":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Orb integration")
			}

			tenantSettings.OrbApiKey = &apiKey

		case "orbit":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Orbit integration")
			}

			tenantSettings.OrbitApiKey = &apiKey

		case "pagerduty":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for PagerDuty integration")
			}

			tenantSettings.PagerDutyApikey = &apiKey

		case "paypaltransaction":
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for PayPal transaction integration")
			}
			secret, ok := data["secret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing secret for PayPal transaction integration")
			}

			tenantSettings.PaypalTransactionClientId = &clientId
			tenantSettings.PaypalTransactionSecret = &secret

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.PaystackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.PaystackLookbackWindow = nil
			}

		case "paystack":
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing secret key for Paystack integration")
			}

			tenantSettings.PaystackSecretKey = &secretKey

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.PaystackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.PaystackLookbackWindow = nil
			}

		case "pendo":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Pendo integration")
			}

			tenantSettings.PendoApiToken = &apiToken

		case "pipedrive":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Pipedrive integration")
			}

			tenantSettings.PipedriveApiToken = &apiToken

		case "plaid":
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing access token for Plaid integration")
			}

			tenantSettings.PlaidAccessToken = &accessToken

		case "plausible":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Plausible integration")
			}
			siteId, ok := data["siteId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing site id for Plausible integration")
			}

			tenantSettings.PlausibleApiKey = &apiKey
			tenantSettings.PlausibleSiteId = &siteId

		case "posthog":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for PostHog integration")
			}

			tenantSettings.PostHogApiKey = &apiKey

			baseUrl, ok := data["baseUrl"].(string)

			if ok && baseUrl != "" {
				tenantSettings.PostHogBaseUrl = &baseUrl
			} else {
				tenantSettings.PostHogBaseUrl = nil
			}

		case "qualaroo":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Qualaroo integration")
			}
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Qualaroo integration")
			}

			tenantSettings.QualarooApiKey = &apiKey
			tenantSettings.QualarooApiToken = &apiToken

		case "quickbooks":
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client id for QuickBooks integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, fmt.Errorf("missing client secret for QuickBooks integration")
			}
			realmId, ok := data["realmId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing realm id for QuickBooks integration")
			}
			refreshToken, ok := data["refreshToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing refresh token for QuickBooks integration")
			}

			tenantSettings.QuickBooksClientId = &clientId
			tenantSettings.QuickBooksClientSecret = &clientSecret
			tenantSettings.QuickBooksRealmId = &realmId
			tenantSettings.QuickBooksRefreshToken = &refreshToken

		case "recharge":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Recharge integration")
			}

			tenantSettings.RechargeApiToken = &apiToken

		case "recruitee":
			companyId, ok := data["companyId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing company id for Recruitee integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Recruitee integration")
			}

			tenantSettings.RecruiteeCompanyId = &companyId
			tenantSettings.RecruiteeApiKey = &apiKey

		case "recurly":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Recurly integration")
			}

			tenantSettings.RecurlyApiKey = &apiKey

		case "retently":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Retently integration")
			}

			tenantSettings.RetentlyApiToken = &apiToken

		case "salesloft":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for SalesLoft integration")
			}

			tenantSettings.SalesloftApiKey = &apiKey

		case "sendgrid":
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API key for Sendgrid integration")
			}

			tenantSettings.SendgridApiKey = &apiKey

		case "slack":
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, fmt.Errorf("missing API token for Slack integration")
			}
			channelFilter, ok := data["channelFilter"].(string)
			if !ok {
				return nil, fmt.Errorf("missing channel filter for Slack integration")
			}

			tenantSettings.SlackApiToken = &apiToken
			tenantSettings.SlackChannelFilter = &channelFilter

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.SlackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.SlackLookbackWindow = nil
			}

		case "stripe":
			accountId, ok := data["accountId"].(string)
			if !ok {
				return nil, fmt.Errorf("missing account id for Stripe integration")
			}
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, fmt.Errorf("missing secret key for Stripe integration")
			}

			tenantSettings.StripeAccountId = &accountId
			tenantSettings.StripeSecretKey = &secretKey

		}

	}

	qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
	if qr.Error != nil {
		return nil, qr.Error
	}
	return qr.Result.(*entity.TenantSettings), nil
}

func (s *tenantSettingsService) ClearIntegrationData(tenantName, identifier string) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		return nil, nil
	} else {

		switch identifier {
		case "hubspot":
			tenantSettings.HubspotPrivateAppKey = nil
		case "zendesk":
			tenantSettings.ZendeskAPIKey = nil
			tenantSettings.ZendeskSubdomain = nil
			tenantSettings.ZendeskAdminEmail = nil
		case "smartsheet":
			tenantSettings.SmartSheetId = nil
			tenantSettings.SmartSheetAccessToken = nil
		case "jira":
			tenantSettings.JiraAPIToken = nil
			tenantSettings.JiraDomain = nil
			tenantSettings.JiraEmail = nil
		case "trello":
			tenantSettings.TrelloAPIToken = nil
			tenantSettings.TrelloAPIKey = nil
		case "aha":
			tenantSettings.AhaAPIUrl = nil
			tenantSettings.AhaAPIKey = nil
		case "airtable":
			tenantSettings.AirtablePersonalAccessToken = nil
		case "amplitude":
			tenantSettings.AmplitudeSecretKey = nil
			tenantSettings.AmplitudeAPIKey = nil
		case "asana":
			tenantSettings.AsanaAccessToken = nil
		case "baton":
			tenantSettings.BatonAPIKey = nil
		case "babelforce":
			tenantSettings.BabelforceRegionEnvironment = nil
			tenantSettings.BabelforceAccessKeyId = nil
			tenantSettings.BabelforceAccessToken = nil
		case "bigquery":
			tenantSettings.BigQueryServiceAccountKey = nil
		case "braintree":
			tenantSettings.BraintreePublicKey = nil
			tenantSettings.BraintreePrivateKey = nil
			tenantSettings.BraintreeEnvironment = nil
			tenantSettings.BraintreeMerchantId = nil
		case "callrail":
			tenantSettings.CallRailAccount = nil
			tenantSettings.CallRailApiToken = nil
		case "chargebee":
			tenantSettings.ChargebeeApiKey = nil
			tenantSettings.ChargebeeProductCatalog = nil
		case "chargify":
			tenantSettings.ChargifyApiKey = nil
			tenantSettings.ChargifyDomain = nil
		case "clickup":
			tenantSettings.ClickUpApiKey = nil
		case "closecom":
			tenantSettings.CloseComApiKey = nil
		case "coda":
			tenantSettings.CodaAuthToken = nil
			tenantSettings.CodaDocumentId = nil
		case "confluence":
			tenantSettings.ConfluenceApiToken = nil
			tenantSettings.ConfluenceDomain = nil
			tenantSettings.ConfluenceLoginEmail = nil
		case "courier":
			tenantSettings.CourierApiKey = nil
		case "customerio":
			tenantSettings.CustomerIoApiKey = nil
		case "datadog":
			tenantSettings.DatadogApiKey = nil
			tenantSettings.DatadogApplicationKey = nil
		case "delighted":
			tenantSettings.DelightedApiKey = nil
		case "dixa":
			tenantSettings.DixaApiToken = nil
		case "drift":
			tenantSettings.DriftApiToken = nil
		case "emailoctopus":
			tenantSettings.EmailOctopusApiKey = nil
		case "facebookMarketing":
			tenantSettings.FacebookMarketingAccessToken = nil
		case "fastbill":
			tenantSettings.FastbillApiKey = nil
			tenantSettings.FastbillProjectId = nil
		case "flexport":
			tenantSettings.FlexportApiKey = nil
		case "freshcaller":
			tenantSettings.FreshcallerApiKey = nil
		case "freshdesk":
			tenantSettings.FreshdeskApiKey = nil
			tenantSettings.FreshdeskDomain = nil
		case "freshsales":
			tenantSettings.FreshsalesApiKey = nil
			tenantSettings.FreshsalesDomain = nil
		case "freshservice":
			tenantSettings.FreshserviceApiKey = nil
			tenantSettings.FreshserviceDomain = nil
		case "genesys":
			tenantSettings.GenesysRegion = nil
			tenantSettings.GenesysClientId = nil
			tenantSettings.GenesysClientSecret = nil
		case "github":
			tenantSettings.GitHubAccessToken = nil
		case "gitlab":
			tenantSettings.GitLabAccessToken = nil
		case "gocardless":
			tenantSettings.GoCardlessAccessToken = nil
			tenantSettings.GoCardlessEnvironment = nil
			tenantSettings.GoCardlessVersion = nil
		case "gong":
			tenantSettings.GongApiKey = nil
		case "harvest":
			tenantSettings.HarvestAccountId = nil
			tenantSettings.HarvestAccessToken = nil
		case "insightly":
			tenantSettings.InsightlyApiToken = nil
		case "instagram":
			tenantSettings.InstagramAccessToken = nil
		case "instatus":
			tenantSettings.InstatusApiKey = nil
		case "intercom":
			tenantSettings.IntercomAccessToken = nil
		case "klaviyo":
			tenantSettings.KlaviyoApiKey = nil
		case "kustomer":
			tenantSettings.KustomerApiToken = nil
		case "looker":
			tenantSettings.LookerClientId = nil
			tenantSettings.LookerClientSecret = nil
			tenantSettings.LookerDomain = nil
		case "mailchimp":
			tenantSettings.MailchimpApiKey = nil
		case "mailjetemail":
			tenantSettings.MailjetEmailApiKey = nil
			tenantSettings.MailjetEmailApiSecret = nil
		case "marketo":
			tenantSettings.MarketoClientId = nil
			tenantSettings.MarketoClientSecret = nil
			tenantSettings.MarketoDomainUrl = nil
		case "microsoftteams":
			tenantSettings.MicrosoftTeamsTenantId = nil
			tenantSettings.MicrosoftTeamsClientId = nil
			tenantSettings.MicrosoftTeamsClientSecret = nil
		case "monday":
			tenantSettings.MondayApiToken = nil
		case "notion":
			tenantSettings.NotionInternalAccessToken = nil
			tenantSettings.NotionPublicClientId = nil
			tenantSettings.NotionPublicClientSecret = nil
			tenantSettings.NotionPublicAccessToken = nil
		case "orb":
			tenantSettings.OrbApiKey = nil
		case "orbit":
			tenantSettings.OrbitApiKey = nil
		case "pagerduty":
			tenantSettings.PagerDutyApikey = nil
		case "paypaltransaction":
			tenantSettings.PaypalTransactionClientId = nil
			tenantSettings.PaypalTransactionSecret = nil
		case "paystack":
			tenantSettings.PaystackSecretKey = nil
			tenantSettings.PaystackLookbackWindow = nil
		case "pendo":
			tenantSettings.PendoApiToken = nil
		case "pipedrive":
			tenantSettings.PipedriveApiToken = nil
		case "plaid":
			tenantSettings.PlaidAccessToken = nil
		case "plausible":
			tenantSettings.PlausibleApiKey = nil
			tenantSettings.PlausibleSiteId = nil
		case "posthog":
			tenantSettings.PostHogApiKey = nil
			tenantSettings.PostHogBaseUrl = nil
		case "qualaroo":
			tenantSettings.QualarooApiKey = nil
			tenantSettings.QualarooApiToken = nil
		case "quickbooks":
			tenantSettings.QuickBooksClientId = nil
			tenantSettings.QuickBooksClientSecret = nil
			tenantSettings.QuickBooksRealmId = nil
			tenantSettings.QuickBooksRefreshToken = nil
		case "recharge":
			tenantSettings.RechargeApiToken = nil
		case "recruitee":
			tenantSettings.RecruiteeCompanyId = nil
			tenantSettings.RecruiteeApiKey = nil
		case "recurly":
			tenantSettings.RecurlyApiKey = nil
		case "retently":
			tenantSettings.RetentlyApiToken = nil
		case "salesloft":
			tenantSettings.SalesloftApiKey = nil
		case "sendgrid":
			tenantSettings.SendgridApiKey = nil
		case "slack":
			tenantSettings.SlackApiToken = nil
			tenantSettings.SlackChannelFilter = nil
			tenantSettings.SlackLookbackWindow = nil
		case "stripe":
			tenantSettings.StripeAccountId = nil
			tenantSettings.StripeSecretKey = nil

		}

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}
