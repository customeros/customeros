package service

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

const GSUITE_SERVICE_PRIVATE_KEY = "GSUITE_SERVICE_PRIVATE_KEY"
const GSUITE_SERVICE_EMAIL_ADDRESS = "GSUITE_SERVICE_EMAIL_ADDRESS"

const SERVICE_GSUITE = "gsuite"
const SERVICE_HUBSPOT = "hubspot"
const SERVICE_SMARTSHEET = "smartsheet"
const SERVICE_JIRA = "jira"
const SERVICE_TRELLO = "trello"
const SERVICE_AHA = "aha"
const SERVICE_AIRTABLE = "airtable"
const SERVICE_AMPLITUDE = "amplitude"
const SERVICE_ASANA = "asana"
const SERVICE_BATON = "baton"
const SERVICE_BABELFORCE = "babelforce"
const SERVICE_BIGQUERY = "bigquery"
const SERVICE_BRAINTREE = "braintree"
const SERVICE_CALLRAIL = "callrail"
const SERVICE_CHARGEBEE = "chargebee"
const SERVICE_CHARGIFY = "chargify"
const SERVICE_CLICKUP = "clickup"
const SERVICE_CLOSECOM = "closecom"
const SERVICE_CODA = "coda"
const SERVICE_CONFLUENCE = "confluence"
const SERVICE_COURIER = "courier"
const SERVICE_CUSTOMERIO = "customerio"
const SERVICE_DATADOG = "datadog"
const SERVICE_DELIGHTED = "delighted"
const SERVICE_DIXA = "dixa"
const SERVICE_DRIFT = "drift"
const SERVICE_EMAILOCTOPUS = "emailoctopus"
const SERVICE_FACEBOOK_MARKETING = "facebookMarketing"
const SERVICE_FASTBILL = "fastbill"
const SERVICE_FLEXPORT = "flexport"
const SERVICE_FRESHCALLER = "freshcaller"
const SERVICE_FRESHDESK = "freshdesk"
const SERVICE_FRESHSALES = "freshsales"
const SERVICE_FRESHSERVICE = "freshservice"
const SERVICE_GENESYS = "genesys"
const SERVICE_GITHUB = "github"
const SERVICE_GITLAB = "gitlab"
const SERVICE_GOCARDLESS = "gocardless"
const SERVICE_GONG = "gong"
const SERVICE_HARVEST = "harvest"
const SERVICE_INSIGHTLY = "insightly"
const SERVICE_INSTAGRAM = "instagram"
const SERVICE_INSTATUS = "instatus"
const SERVICE_INTERCOM = "intercom"
const SERVICE_KLAVIYO = "klaviyo"
const SERVICE_KUSTOMER = "kustomer"
const SERVICE_LOOKER = "looker"
const SERVICE_MAILCHIMP = "mailchimp"
const SERVICE_MAILJETEMAIL = "mailjetemail"
const SERVICE_MARKETO = "marketo"
const SERVICE_MICROSOFT_TEAMS = "microsoftteams"
const SERVICE_MONDAY = "monday"
const SERVICE_NOTION = "notion"
const SERVICE_ORB = "orb"
const SERVICE_ORACLE_NETSUITE = "oraclenetsuite"
const SERVICE_ORBIT = "orbit"
const SERVICE_PAGERDUTY = "pagerduty"
const SERVICE_PAYSTACK = "paystack"
const SERVICE_PENDO = "pendo"
const SERVICE_PIPEDRIVE = "pipedrive"
const SERVICE_PLAID = "plaid"
const SERVICE_PLAUSIBLE = "plausible"
const SERVICE_PAYPAL_TRANSACTION = "paypaltransaction"
const SERVICE_POSTHOG = "posthog"
const SERVICE_QUALAROO = "qualaroo"
const SERVICE_QUICKBOOKS = "quickbooks"
const SERVICE_RECHARGE = "recharge"
const SERVICE_RECRUITEE = "recruitee"
const SERVICE_RECURLY = "recurly"
const SERVICE_RETENTLY = "retently"
const SERVICE_SALESFORCE = "salesforce"
const SERVICE_SALESLOFT = "salesloft"
const SERVICE_SENDGRID = "sendgrid"
const SERVICE_SENTRY = "sentry"
const SERVICE_SLACK = "slack"
const SERVICE_STRIPE = "stripe"
const SERVICE_SURVEYSPARROW = "surveysparrow"
const SERVICE_SURVEYMONKEY = "surveymonkey"
const SERVICE_TALKDESK = "talkdesk"
const SERVICE_TIKTOK = "tiktok"
const SERVICE_TODOIST = "todoist"
const SERVICE_TYPEFORM = "typeform"
const SERVICE_VITTALLY = "vittally"
const SERVICE_WRIKE = "wrike"
const SERVICE_XERO = "xero"
const SERVICE_ZENDESK_SUPPORT = "zendesksupport"
const SERVICE_ZENDESK_CHAT = "zendeskchat"
const SERVICE_ZENDESK_TALK = "zendesktalk"
const SERVICE_ZENDESK_SELL = "zendesksell"
const SERVICE_ZENDESK_SUNSHINE = "zendesksunshine"
const SERVICE_ZENEFITS = "zenefits"
const SERVICE_MIXPANEL = "mixpanel"

type TenantSettingsService interface {
	GetForTenant(tenantName string) (*entity.TenantSettings, map[string]bool, error)
	SaveIntegrationData(tenantName string, request map[string]interface{}) (*entity.TenantSettings, map[string]bool, error)
	ClearIntegrationData(tenantName, identifier string) (*entity.TenantSettings, map[string]bool, error)
}

type tenantSettingsService struct {
	repositories *repository.PostgresRepositories
	serviceMap   map[string][]keyMapping
	log          logger.Logger
}

type keyMapping struct {
	ApiKeyName string
	DbKeyName  string
}

func NewTenantSettingsService(repositories *repository.PostgresRepositories, log logger.Logger) TenantSettingsService {
	return &tenantSettingsService{
		repositories: repositories,
		serviceMap: map[string][]keyMapping{
			SERVICE_GSUITE: {
				keyMapping{"privateKey", GSUITE_SERVICE_PRIVATE_KEY},
				keyMapping{"clientEmail", GSUITE_SERVICE_EMAIL_ADDRESS},
			},
		},
		log: log,
	}
}

func (s *tenantSettingsService) GetForTenant(tenantName string) (*entity.TenantSettings, map[string]bool, error) {
	qr := s.repositories.TenantSettingsRepository.FindForTenantName(tenantName)
	var settings entity.TenantSettings
	var ok bool
	if qr.Error != nil {
		return nil, nil, qr.Error
	} else if qr.Result == nil {
		return nil, nil, nil
	} else {
		settings, ok = qr.Result.(entity.TenantSettings)
		if !ok {
			return nil, nil, fmt.Errorf("GetForTenant: unexpected type %T", qr.Result)
		}
	}
	activeServices, err := s.GetServiceActivations(tenantName)
	if err != nil {
		return nil, nil, fmt.Errorf("SaveIntegrationData: %v", err)
	}
	return &settings, activeServices, nil
}

func (s *tenantSettingsService) GetServiceActivations(tenantName string) (map[string]bool, error) {
	result := make(map[string]bool)
	for service, keyMappings := range s.serviceMap {
		keys := make([]string, 0)
		for _, mapping := range keyMappings {
			keys = append(keys, mapping.DbKeyName)
		}
		active, err := s.repositories.TenantSettingsRepository.CheckKeysExist(tenantName, keys)
		if err != nil {
			return nil, fmt.Errorf("GetServiceActivations: %w", err)
		}
		result[service] = active
	}

	return result, nil
}

func (s *tenantSettingsService) SaveIntegrationData(tenantName string, request map[string]interface{}) (*entity.TenantSettings, map[string]bool, error) {
	tenantSettings, _, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, nil, err
	}
	var keysToUpdate []entity.TenantAPIKey
	legacyUpdate := false

	if tenantSettings == nil {
		tenantSettings = &entity.TenantSettings{
			TenantName: tenantName,
		}

		if qr := s.repositories.TenantSettingsRepository.Save(tenantSettings); qr.Error != nil {
			return nil, nil, qr.Error
		}
	}

	// Update tenant settings with new integration data
	for integrationId, value := range request {
		data, ok := value.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("invalid data for integration %s", integrationId)
		}

		mappings, ok := s.serviceMap[integrationId]
		if ok {
			for _, mapping := range mappings {
				if value, ok := data[mapping.ApiKeyName]; ok {
					valueStr, ok := value.(string)
					if !ok {
						return nil, nil, fmt.Errorf("invalid data for key %s in integration %s", mapping.ApiKeyName, integrationId)
					}
					keysToUpdate = append(keysToUpdate, entity.TenantAPIKey{TenantName: tenantName, Key: mapping.DbKeyName, Value: valueStr})
					data[mapping.DbKeyName] = value
				}
			}
			continue
		}

		legacyUpdate = true

		switch integrationId {
		case SERVICE_HUBSPOT:
			privateAppKey, ok := data["privateAppKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing private app key for Hubspot integration")
			}
			tenantSettings.HubspotPrivateAppKey = &privateAppKey

		case SERVICE_SMARTSHEET:
			id, ok := data["id"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing Smartsheet ID")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Smartsheet integration")
			}
			tenantSettings.SmartSheetId = &id
			tenantSettings.SmartSheetAccessToken = &accessToken

		case SERVICE_JIRA:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Jira integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Jira integration")
			}
			email, ok := data["email"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing email for Jira integration")
			}
			tenantSettings.JiraAPIToken = &apiToken
			tenantSettings.JiraDomain = &domain
			tenantSettings.JiraEmail = &email

		case SERVICE_TRELLO:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Trello integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Trello integration")
			}
			tenantSettings.TrelloAPIToken = &apiToken
			tenantSettings.TrelloAPIKey = &apiKey

		case SERVICE_AHA:
			apiUrl, ok := data["apiUrl"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API Url for Aha integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Aha integration")
			}
			tenantSettings.AhaAPIUrl = &apiUrl
			tenantSettings.AhaAPIKey = &apiKey

		case SERVICE_AIRTABLE:
			personalAccessToken, ok := data["personalAccessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing personal access token for Airtable integration")
			}
			tenantSettings.AirtablePersonalAccessToken = &personalAccessToken

		case SERVICE_AMPLITUDE:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Amplitude integration")
			}
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing secret key for Amplitude integration")
			}
			tenantSettings.AmplitudeSecretKey = &secretKey
			tenantSettings.AmplitudeAPIKey = &apiKey

		case SERVICE_ASANA:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Asana integration")
			}

			tenantSettings.AsanaAccessToken = &accessToken

		case SERVICE_BATON:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Baton integration")
			}
			tenantSettings.BatonAPIKey = &apiKey

		case SERVICE_BABELFORCE:
			regionEnvironment, ok := data["regionEnvironment"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing region / environment for Babelforce integration")
			}
			accessKeyId, ok := data["accessKeyId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access key id for Babelforce integration")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Babelforce integration")
			}

			tenantSettings.BabelforceRegionEnvironment = &regionEnvironment
			tenantSettings.BabelforceAccessKeyId = &accessKeyId
			tenantSettings.BabelforceAccessToken = &accessToken

		case SERVICE_BIGQUERY:
			serviceAccountKey, ok := data["serviceAccountKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing service account key for BigQuery integration")
			}

			tenantSettings.BigQueryServiceAccountKey = &serviceAccountKey

		case SERVICE_BRAINTREE:
			publicKey, ok := data["publicKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing public key for Braintree integration")
			}
			privateKey, ok := data["privateKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing private key for Braintree integration")
			}
			environment, ok := data["environment"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing environment for Braintree integration")
			}
			merchantId, ok := data["merchantId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing merchant id for Braintree integration")
			}

			tenantSettings.BraintreePublicKey = &publicKey
			tenantSettings.BraintreePrivateKey = &privateKey
			tenantSettings.BraintreeEnvironment = &environment
			tenantSettings.BraintreeMerchantId = &merchantId

		case SERVICE_CALLRAIL:
			account, ok := data["account"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing account for CallRail integration")
			}
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for CallRail integration")
			}

			tenantSettings.CallRailAccount = &account
			tenantSettings.CallRailApiToken = &apiToken

		case SERVICE_CHARGEBEE:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Chargebee integration")
			}
			productCatalog, ok := data["productCatalog"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing product catalog for CallRail integration")
			}

			tenantSettings.ChargebeeApiKey = &apiKey
			tenantSettings.ChargebeeProductCatalog = &productCatalog

		case SERVICE_CHARGIFY:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Chargify integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Chargify integration")
			}

			tenantSettings.ChargifyApiKey = &apiKey
			tenantSettings.ChargifyDomain = &domain

		case SERVICE_CLICKUP:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for ClickUp integration")
			}

			tenantSettings.ClickUpApiKey = &apiKey

		case SERVICE_CLOSECOM:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Close.com integration")
			}

			tenantSettings.CloseComApiKey = &apiKey

		case SERVICE_CODA:
			authToken, ok := data["authToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing auth token for Coda integration")
			}
			documentId, ok := data["documentId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing document id for Coda integration")
			}

			tenantSettings.CodaAuthToken = &authToken
			tenantSettings.CodaDocumentId = &documentId

		case SERVICE_CONFLUENCE:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Confluence integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Confluence integration")
			}
			loginEmail, ok := data["loginEmail"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing login email for Confluence integration")
			}

			tenantSettings.ConfluenceApiToken = &apiToken
			tenantSettings.ConfluenceDomain = &domain
			tenantSettings.ConfluenceLoginEmail = &loginEmail

		case SERVICE_COURIER:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Courier integration")
			}

			tenantSettings.CourierApiKey = &apiKey

		case SERVICE_CUSTOMERIO:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Customer.io integration")
			}

			tenantSettings.CustomerIoApiKey = &apiKey

		case SERVICE_DATADOG:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Customer.io integration")
			}
			applicationKey, ok := data["applicationKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing application key for Customer.io integration")
			}

			tenantSettings.DatadogApiKey = &apiKey
			tenantSettings.DatadogApplicationKey = &applicationKey

		case SERVICE_DELIGHTED:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Delighted integration")
			}

			tenantSettings.DelightedApiKey = &apiKey

		case SERVICE_DIXA:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Dixa integration")
			}

			tenantSettings.DixaApiToken = &apiToken

		case SERVICE_DRIFT:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Drift integration")
			}

			tenantSettings.DriftApiToken = &apiToken

		case SERVICE_EMAILOCTOPUS:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for EmailOctopus integration")
			}

			tenantSettings.EmailOctopusApiKey = &apiKey

		case SERVICE_FACEBOOK_MARKETING:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Facebook integration")
			}

			tenantSettings.FacebookMarketingAccessToken = &accessToken

		case SERVICE_FASTBILL:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Fastbill integration")
			}
			projectId, ok := data["projectId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing project id for Fastbill integration")
			}

			tenantSettings.FastbillApiKey = &apiKey
			tenantSettings.FastbillProjectId = &projectId

		case SERVICE_FLEXPORT:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Flexport integration")
			}

			tenantSettings.FlexportApiKey = &apiKey

		case SERVICE_FRESHCALLER:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Freshcaller integration")
			}

			tenantSettings.FreshcallerApiKey = &apiKey

		case SERVICE_FRESHDESK:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Freshdesk integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Freshdesk integration")
			}

			tenantSettings.FreshdeskApiKey = &apiKey
			tenantSettings.FreshdeskDomain = &domain

		case SERVICE_FRESHSALES:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Freshsales integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Freshsales integration")
			}

			tenantSettings.FreshsalesApiKey = &apiKey
			tenantSettings.FreshsalesDomain = &domain

		case SERVICE_FRESHSERVICE:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Freshservice integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Freshservice integration")
			}

			tenantSettings.FreshserviceApiKey = &apiKey
			tenantSettings.FreshserviceDomain = &domain

		case SERVICE_GENESYS:
			region, ok := data["region"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing region for Genesys integration")
			}
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Genesys integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Genesys integration")
			}

			tenantSettings.GenesysRegion = &region
			tenantSettings.GenesysClientId = &clientId
			tenantSettings.GenesysClientSecret = &clientSecret

		case SERVICE_GITHUB:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for GitHub integration")
			}

			tenantSettings.GitHubAccessToken = &accessToken

		case SERVICE_GITLAB:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for GitLab integration")
			}

			tenantSettings.GitLabAccessToken = &accessToken

		case SERVICE_GOCARDLESS:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for GoCardless integration")
			}
			environment, ok := data["environment"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing environment for GoCardless integration")
			}
			version, ok := data["version"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing version for GoCardless integration")
			}

			tenantSettings.GoCardlessAccessToken = &accessToken
			tenantSettings.GoCardlessEnvironment = &environment
			tenantSettings.GoCardlessVersion = &version

		case SERVICE_GONG:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Gong integration")
			}

			tenantSettings.GongApiKey = &apiKey

		case SERVICE_HARVEST:
			accountId, ok := data["accountId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing account id for Harvest integration")
			}
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Harvest integration")
			}

			tenantSettings.HarvestAccountId = &accountId
			tenantSettings.HarvestAccessToken = &accessToken

		case SERVICE_INSIGHTLY:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Insightly integration")
			}

			tenantSettings.InsightlyApiToken = &apiToken

		case SERVICE_INSTAGRAM:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Harvest integration")
			}

			tenantSettings.InstagramAccessToken = &accessToken

		case SERVICE_INSTATUS:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Instatus integration")
			}

			tenantSettings.InstatusApiKey = &apiKey

		case SERVICE_INTERCOM:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Intercom integration")
			}

			tenantSettings.IntercomAccessToken = &accessToken

		case SERVICE_KLAVIYO:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Klaviyo integration")
			}

			tenantSettings.KlaviyoApiKey = &apiKey

		case SERVICE_KUSTOMER:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Kustomer integration")
			}

			tenantSettings.KustomerApiToken = &apiToken

		case SERVICE_LOOKER:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Looker integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Looker integration")
			}
			domain, ok := data["domain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain for Looker integration")
			}

			tenantSettings.LookerClientId = &clientId
			tenantSettings.LookerClientSecret = &clientSecret
			tenantSettings.LookerDomain = &domain

		case SERVICE_MAILCHIMP:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Mailchimp integration")
			}

			tenantSettings.MailchimpApiKey = &apiKey

		case SERVICE_MAILJETEMAIL:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Mailjet Email integration")
			}
			apiSecret, ok := data["apiSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API secret for Mailjet Email integration")
			}

			tenantSettings.MailjetEmailApiKey = &apiKey
			tenantSettings.MailjetEmailApiSecret = &apiSecret

		case SERVICE_MARKETO:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Marketo integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Marketo integration")
			}
			domainUrl, ok := data["domainUrl"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing domain URL for Marketo integration")
			}

			tenantSettings.MarketoClientId = &clientId
			tenantSettings.MarketoClientSecret = &clientSecret
			tenantSettings.MarketoDomainUrl = &domainUrl

		case SERVICE_MICROSOFT_TEAMS:
			tenantId, ok := data["tenantId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing tenant id for Microsoft Teams integration")
			}
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Microsoft Teams integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Microsoft Teams integration")
			}

			tenantSettings.MicrosoftTeamsTenantId = &tenantId
			tenantSettings.MicrosoftTeamsClientId = &clientId
			tenantSettings.MicrosoftTeamsClientSecret = &clientSecret

		case SERVICE_MONDAY:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Monday integration")
			}

			tenantSettings.MondayApiToken = &apiToken

		case SERVICE_NOTION:
			internalAccessToken, _ := data["internalAccessToken"].(string)

			publicClientId, _ := data["publicClientId"].(string)
			publicClientSecret, _ := data["publicClientSecret"].(string)
			publicAccessToken, _ := data["publicAccessToken"].(string)

			if internalAccessToken == "" && publicClientId == "" && publicClientSecret == "" && publicAccessToken == "" {
				return nil, nil, fmt.Errorf("missing Notion integration data")
			}

			if internalAccessToken == "" && (publicClientId == "" || publicClientSecret == "" || publicAccessToken == "") {
				return nil, nil, fmt.Errorf("missing public Notion integration data")
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

		case SERVICE_ORACLE_NETSUITE:
			accountId, ok := data["accountId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing account id for Oracle Netsuite integration")
			}
			consumerKey, ok := data["consumerKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing consumer key for Oracle Netsuite integration")
			}
			consumerSecret, ok := data["consumerSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing consumer secret for Oracle Netsuite integration")
			}
			tokenId, ok := data["tokenId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing token id for Oracle Netsuite integration")
			}
			tokenSecret, ok := data["tokenSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing token secret for Oracle Netsuite integration")
			}

			tenantSettings.OracleNetsuiteAccountId = &accountId
			tenantSettings.OracleNetsuiteConsumerKey = &consumerKey
			tenantSettings.OracleNetsuiteConsumerSecret = &consumerSecret
			tenantSettings.OracleNetsuiteTokenId = &tokenId
			tenantSettings.OracleNetsuiteTokenSecret = &tokenSecret

		case SERVICE_ORB:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Orb integration")
			}

			tenantSettings.OrbApiKey = &apiKey

		case SERVICE_ORBIT:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Orbit integration")
			}

			tenantSettings.OrbitApiKey = &apiKey

		case SERVICE_PAGERDUTY:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for PagerDuty integration")
			}

			tenantSettings.PagerDutyApikey = &apiKey

		case SERVICE_PAYPAL_TRANSACTION:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for PayPal transaction integration")
			}
			secret, ok := data["secret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing secret for PayPal transaction integration")
			}

			tenantSettings.PaypalTransactionClientId = &clientId
			tenantSettings.PaypalTransactionSecret = &secret

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.PaystackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.PaystackLookbackWindow = nil
			}

		case SERVICE_PAYSTACK:
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing secret key for Paystack integration")
			}

			tenantSettings.PaystackSecretKey = &secretKey

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.PaystackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.PaystackLookbackWindow = nil
			}

		case SERVICE_PENDO:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Pendo integration")
			}

			tenantSettings.PendoApiToken = &apiToken

		case SERVICE_PIPEDRIVE:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Pipedrive integration")
			}

			tenantSettings.PipedriveApiToken = &apiToken

		case SERVICE_PLAID:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Plaid integration")
			}

			tenantSettings.PlaidAccessToken = &accessToken

		case SERVICE_PLAUSIBLE:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Plausible integration")
			}
			siteId, ok := data["siteId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing site id for Plausible integration")
			}

			tenantSettings.PlausibleApiKey = &apiKey
			tenantSettings.PlausibleSiteId = &siteId

		case SERVICE_POSTHOG:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for PostHog integration")
			}

			tenantSettings.PostHogApiKey = &apiKey

			baseUrl, ok := data["baseUrl"].(string)

			if ok && baseUrl != "" {
				tenantSettings.PostHogBaseUrl = &baseUrl
			} else {
				tenantSettings.PostHogBaseUrl = nil
			}

		case SERVICE_QUALAROO:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Qualaroo integration")
			}
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Qualaroo integration")
			}

			tenantSettings.QualarooApiKey = &apiKey
			tenantSettings.QualarooApiToken = &apiToken

		case SERVICE_QUICKBOOKS:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for QuickBooks integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for QuickBooks integration")
			}
			realmId, ok := data["realmId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing realm id for QuickBooks integration")
			}
			refreshToken, ok := data["refreshToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing refresh token for QuickBooks integration")
			}

			tenantSettings.QuickBooksClientId = &clientId
			tenantSettings.QuickBooksClientSecret = &clientSecret
			tenantSettings.QuickBooksRealmId = &realmId
			tenantSettings.QuickBooksRefreshToken = &refreshToken

		case SERVICE_RECHARGE:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Recharge integration")
			}

			tenantSettings.RechargeApiToken = &apiToken

		case SERVICE_RECRUITEE:
			companyId, ok := data["companyId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing company id for Recruitee integration")
			}
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Recruitee integration")
			}

			tenantSettings.RecruiteeCompanyId = &companyId
			tenantSettings.RecruiteeApiKey = &apiKey

		case SERVICE_RECURLY:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Recurly integration")
			}

			tenantSettings.RecurlyApiKey = &apiKey

		case SERVICE_RETENTLY:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Retently integration")
			}

			tenantSettings.RetentlyApiToken = &apiToken

		case SERVICE_SALESFORCE:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Salesforce integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Salesforce integration")
			}
			refreshToken, ok := data["refreshToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing refresh token for Salesforce integration")
			}

			tenantSettings.SalesforceClientId = &clientId
			tenantSettings.SalesforceClientSecret = &clientSecret
			tenantSettings.SalesforceRefreshToken = &refreshToken

		case SERVICE_SALESLOFT:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for SalesLoft integration")
			}

			tenantSettings.SalesloftApiKey = &apiKey

		case SERVICE_SENDGRID:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Sendgrid integration")
			}

			tenantSettings.SendgridApiKey = &apiKey

		case SERVICE_SENTRY:
			project, ok := data["project"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing project for Sentry integration")
			}
			authenticationToken, ok := data["authenticationToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing authentication token for Sentry integration")
			}
			organization, ok := data["organization"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing organization for Sentry integration")
			}

			tenantSettings.SentryProject = &project
			tenantSettings.SentryAuthenticationToken = &authenticationToken
			tenantSettings.SentryOrganization = &organization

			host, ok := data["host"].(string)

			if ok && host != "" {
				tenantSettings.SentryHost = &host
			} else {
				tenantSettings.SentryHost = nil
			}

		case SERVICE_SLACK:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Slack integration")
			}
			channelFilter, ok := data["channelFilter"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing channel filter for Slack integration")
			}

			tenantSettings.SlackApiToken = &apiToken
			tenantSettings.SlackChannelFilter = &channelFilter

			lookbackWindow, ok := data["lookbackWindow"].(string)

			if ok && lookbackWindow != "" {
				tenantSettings.SlackLookbackWindow = &lookbackWindow
			} else {
				tenantSettings.SlackLookbackWindow = nil
			}

		case SERVICE_STRIPE:
			accountId, ok := data["accountId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing account id for Stripe integration")
			}
			secretKey, ok := data["secretKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing secret key for Stripe integration")
			}

			tenantSettings.StripeAccountId = &accountId
			tenantSettings.StripeSecretKey = &secretKey

		case SERVICE_SURVEYSPARROW:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for SurveySparrow integration")
			}

			tenantSettings.SurveySparrowAccessToken = &accessToken

		case SERVICE_SURVEYMONKEY:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for SurveyMonkey integration")
			}

			tenantSettings.SurveyMonkeyAccessToken = &accessToken

		case SERVICE_TALKDESK:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Talkdesk integration")
			}

			tenantSettings.TalkdeskApiKey = &apiKey

		case SERVICE_TIKTOK:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for TikTok integration")
			}

			tenantSettings.TikTokAccessToken = &accessToken

		case SERVICE_TODOIST:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Todoist integration")
			}

			tenantSettings.TodoistApiToken = &apiToken

		case SERVICE_TYPEFORM:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Typeform integration")
			}

			tenantSettings.TypeformApiToken = &apiToken

		case SERVICE_VITTALLY:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Vittally integration")
			}

			tenantSettings.VittallyApiKey = &apiKey

		case SERVICE_WRIKE:
			accessToken, ok := data["accessToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access token for Wrike integration")
			}
			hostUrl, ok := data["hostUrl"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing host url for Wrike integration")
			}

			tenantSettings.WrikeAccessToken = &accessToken
			tenantSettings.WrikeHostUrl = &hostUrl

		case SERVICE_XERO:
			clientId, ok := data["clientId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Xero integration")
			}
			clientSecret, ok := data["clientSecret"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client secret for Xero integration")
			}
			tenantId, ok := data["tenantId"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing tenant id for Xero integration")
			}
			scopes, ok := data["scopes"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing scopes for Xero integration")
			}

			tenantSettings.XeroClientId = &clientId
			tenantSettings.XeroClientSecret = &clientSecret
			tenantSettings.XeroTenantId = &tenantId
			tenantSettings.XeroScopes = &scopes

		case SERVICE_ZENDESK_SUPPORT:
			apiKey, ok := data["apiKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API key for Zendesk integration")
			}
			subdomain, ok := data["subdomain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing subdomain for Zendesk integration")
			}
			adminEmail, ok := data["adminEmail"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing admin email for Zendesk integration")
			}
			tenantSettings.ZendeskAPIKey = &apiKey
			tenantSettings.ZendeskSubdomain = &subdomain
			tenantSettings.ZendeskAdminEmail = &adminEmail

		case SERVICE_ZENDESK_CHAT:
			subdomain, ok := data["subdomain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing subdomain for Zendesk Chat integration")
			}
			accessKey, ok := data["accessKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access key for Zendesk Chat integration")
			}

			tenantSettings.ZendeskChatSubdomain = &subdomain
			tenantSettings.ZendeskChatAccessKey = &accessKey

		case SERVICE_ZENDESK_TALK:
			subdomain, ok := data["subdomain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing subdomain for Zendesk Talk integration")
			}
			accessKey, ok := data["accessKey"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing access key for Zendesk Talk integration")
			}

			tenantSettings.ZendeskTalkSubdomain = &subdomain
			tenantSettings.ZendeskTalkAccessKey = &accessKey

		case SERVICE_ZENDESK_SELL:
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Zendesk Sell integration")
			}

			tenantSettings.ZendeskSellApiToken = &apiToken

		case SERVICE_ZENDESK_SUNSHINE:
			subdomain, ok := data["subdomain"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing subdomain for Zendesk Sunshine integration")
			}
			apiToken, ok := data["apiToken"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing API token for Zendesk Sunshine integration")
			}
			email, ok := data["email"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing email for Zendesk Sunshine integration")
			}

			tenantSettings.ZendeskSunshineSubdomain = &subdomain
			tenantSettings.ZendeskSunshineApiToken = &apiToken
			tenantSettings.ZendeskSunshineEmail = &email

		case SERVICE_ZENEFITS:
			token, ok := data["token"].(string)
			if !ok {
				return nil, nil, fmt.Errorf("missing client id for Xero integration")
			}

			tenantSettings.ZenefitsToken = &token

		case SERVICE_MIXPANEL:
			username, _ := data["username"].(string)
			secret, _ := data["secret"].(string)
			projectId, _ := data["projectId"].(string)
			projectSecret, _ := data["projectSecret"].(string)
			projectTimezone, _ := data["projectTimezone"].(string)
			region, _ := data["region"].(string)

			tenantSettings.MixpanelUsername = &username
			tenantSettings.MixpanelSecret = &secret
			tenantSettings.MixpanelProjectId = &projectId
			tenantSettings.MixpanelProjectSecret = &projectSecret
			tenantSettings.MixpanelProjectTimezone = &projectTimezone
			tenantSettings.MixpanelRegion = &region
		}
	}

	if legacyUpdate {
		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, nil, fmt.Errorf("SaveIntegrationData: %v", qr.Error)
		}
		tenantSettings = qr.Result.(*entity.TenantSettings)
	}

	if keysToUpdate != nil {
		err = s.repositories.TenantSettingsRepository.SaveKeys(keysToUpdate)
		if err != nil {
			return nil, nil, fmt.Errorf("SaveIntegrationData: %v", err)
		}
	}

	activeServices, err := s.GetServiceActivations(tenantName)
	if err != nil {
		return nil, nil, fmt.Errorf("SaveIntegrationData: %v", err)
	}

	return tenantSettings, activeServices, nil
}

func (s *tenantSettingsService) ClearIntegrationData(tenantName, identifier string) (*entity.TenantSettings, map[string]bool, error) {
	tenantSettings, _, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, nil, err
	}

	if tenantSettings == nil {
		return nil, nil, nil
	} else {

		var keysToDelete []entity.TenantAPIKey
		mappings, ok := s.serviceMap[identifier]
		if ok {
			for _, mapping := range mappings {

				keysToDelete = append(keysToDelete, entity.TenantAPIKey{TenantName: tenantName, Key: mapping.DbKeyName})
			}
			err = s.repositories.TenantSettingsRepository.DeleteKeys(keysToDelete)
			if err != nil {
				return nil, nil, fmt.Errorf("ClearIntegrationData: %v", err)
			}
		} else {

			switch identifier {
			case SERVICE_HUBSPOT:
				tenantSettings.HubspotPrivateAppKey = nil
			case SERVICE_SMARTSHEET:
				tenantSettings.SmartSheetId = nil
				tenantSettings.SmartSheetAccessToken = nil
			case SERVICE_JIRA:
				tenantSettings.JiraAPIToken = nil
				tenantSettings.JiraDomain = nil
				tenantSettings.JiraEmail = nil
			case SERVICE_TRELLO:
				tenantSettings.TrelloAPIToken = nil
				tenantSettings.TrelloAPIKey = nil
			case SERVICE_AHA:
				tenantSettings.AhaAPIUrl = nil
				tenantSettings.AhaAPIKey = nil
			case SERVICE_AIRTABLE:
				tenantSettings.AirtablePersonalAccessToken = nil
			case SERVICE_AMPLITUDE:
				tenantSettings.AmplitudeSecretKey = nil
				tenantSettings.AmplitudeAPIKey = nil
			case SERVICE_ASANA:
				tenantSettings.AsanaAccessToken = nil
			case SERVICE_BATON:
				tenantSettings.BatonAPIKey = nil
			case SERVICE_BABELFORCE:
				tenantSettings.BabelforceRegionEnvironment = nil
				tenantSettings.BabelforceAccessKeyId = nil
				tenantSettings.BabelforceAccessToken = nil
			case SERVICE_BIGQUERY:
				tenantSettings.BigQueryServiceAccountKey = nil
			case SERVICE_BRAINTREE:
				tenantSettings.BraintreePublicKey = nil
				tenantSettings.BraintreePrivateKey = nil
				tenantSettings.BraintreeEnvironment = nil
				tenantSettings.BraintreeMerchantId = nil
			case SERVICE_CALLRAIL:
				tenantSettings.CallRailAccount = nil
				tenantSettings.CallRailApiToken = nil
			case SERVICE_CHARGEBEE:
				tenantSettings.ChargebeeApiKey = nil
				tenantSettings.ChargebeeProductCatalog = nil
			case SERVICE_CHARGIFY:
				tenantSettings.ChargifyApiKey = nil
				tenantSettings.ChargifyDomain = nil
			case SERVICE_CLICKUP:
				tenantSettings.ClickUpApiKey = nil
			case SERVICE_CLOSECOM:
				tenantSettings.CloseComApiKey = nil
			case SERVICE_CODA:
				tenantSettings.CodaAuthToken = nil
				tenantSettings.CodaDocumentId = nil
			case SERVICE_CONFLUENCE:
				tenantSettings.ConfluenceApiToken = nil
				tenantSettings.ConfluenceDomain = nil
				tenantSettings.ConfluenceLoginEmail = nil
			case SERVICE_COURIER:
				tenantSettings.CourierApiKey = nil
			case SERVICE_CUSTOMERIO:
				tenantSettings.CustomerIoApiKey = nil
			case SERVICE_DATADOG:
				tenantSettings.DatadogApiKey = nil
				tenantSettings.DatadogApplicationKey = nil
			case SERVICE_DELIGHTED:
				tenantSettings.DelightedApiKey = nil
			case SERVICE_DIXA:
				tenantSettings.DixaApiToken = nil
			case SERVICE_DRIFT:
				tenantSettings.DriftApiToken = nil
			case SERVICE_EMAILOCTOPUS:
				tenantSettings.EmailOctopusApiKey = nil
			case SERVICE_FACEBOOK_MARKETING:
				tenantSettings.FacebookMarketingAccessToken = nil
			case SERVICE_FASTBILL:
				tenantSettings.FastbillApiKey = nil
				tenantSettings.FastbillProjectId = nil
			case SERVICE_FLEXPORT:
				tenantSettings.FlexportApiKey = nil
			case SERVICE_FRESHCALLER:
				tenantSettings.FreshcallerApiKey = nil
			case SERVICE_FRESHDESK:
				tenantSettings.FreshdeskApiKey = nil
				tenantSettings.FreshdeskDomain = nil
			case SERVICE_FRESHSALES:
				tenantSettings.FreshsalesApiKey = nil
				tenantSettings.FreshsalesDomain = nil
			case SERVICE_FRESHSERVICE:
				tenantSettings.FreshserviceApiKey = nil
				tenantSettings.FreshserviceDomain = nil
			case SERVICE_GENESYS:
				tenantSettings.GenesysRegion = nil
				tenantSettings.GenesysClientId = nil
				tenantSettings.GenesysClientSecret = nil
			case SERVICE_GITHUB:
				tenantSettings.GitHubAccessToken = nil
			case SERVICE_GITLAB:
				tenantSettings.GitLabAccessToken = nil
			case SERVICE_GOCARDLESS:
				tenantSettings.GoCardlessAccessToken = nil
				tenantSettings.GoCardlessEnvironment = nil
				tenantSettings.GoCardlessVersion = nil
			case SERVICE_GONG:
				tenantSettings.GongApiKey = nil
			case SERVICE_HARVEST:
				tenantSettings.HarvestAccountId = nil
				tenantSettings.HarvestAccessToken = nil
			case SERVICE_INSIGHTLY:
				tenantSettings.InsightlyApiToken = nil
			case SERVICE_INSTAGRAM:
				tenantSettings.InstagramAccessToken = nil
			case SERVICE_INSTATUS:
				tenantSettings.InstatusApiKey = nil
			case SERVICE_INTERCOM:
				tenantSettings.IntercomAccessToken = nil
			case SERVICE_KLAVIYO:
				tenantSettings.KlaviyoApiKey = nil
			case SERVICE_KUSTOMER:
				tenantSettings.KustomerApiToken = nil
			case SERVICE_LOOKER:
				tenantSettings.LookerClientId = nil
				tenantSettings.LookerClientSecret = nil
				tenantSettings.LookerDomain = nil
			case SERVICE_MAILCHIMP:
				tenantSettings.MailchimpApiKey = nil
			case SERVICE_MAILJETEMAIL:
				tenantSettings.MailjetEmailApiKey = nil
				tenantSettings.MailjetEmailApiSecret = nil
			case SERVICE_MARKETO:
				tenantSettings.MarketoClientId = nil
				tenantSettings.MarketoClientSecret = nil
				tenantSettings.MarketoDomainUrl = nil
			case SERVICE_MICROSOFT_TEAMS:
				tenantSettings.MicrosoftTeamsTenantId = nil
				tenantSettings.MicrosoftTeamsClientId = nil
				tenantSettings.MicrosoftTeamsClientSecret = nil
			case SERVICE_MONDAY:
				tenantSettings.MondayApiToken = nil
			case SERVICE_NOTION:
				tenantSettings.NotionInternalAccessToken = nil
				tenantSettings.NotionPublicClientId = nil
				tenantSettings.NotionPublicClientSecret = nil
				tenantSettings.NotionPublicAccessToken = nil
			case SERVICE_ORACLE_NETSUITE:
				tenantSettings.OracleNetsuiteAccountId = nil
				tenantSettings.OracleNetsuiteConsumerKey = nil
				tenantSettings.OracleNetsuiteConsumerSecret = nil
				tenantSettings.OracleNetsuiteTokenId = nil
				tenantSettings.OracleNetsuiteTokenSecret = nil
			case SERVICE_ORB:
				tenantSettings.OrbApiKey = nil
			case SERVICE_ORBIT:
				tenantSettings.OrbitApiKey = nil
			case SERVICE_PAGERDUTY:
				tenantSettings.PagerDutyApikey = nil
			case SERVICE_PAYPAL_TRANSACTION:
				tenantSettings.PaypalTransactionClientId = nil
				tenantSettings.PaypalTransactionSecret = nil
			case SERVICE_PAYSTACK:
				tenantSettings.PaystackSecretKey = nil
				tenantSettings.PaystackLookbackWindow = nil
			case SERVICE_PENDO:
				tenantSettings.PendoApiToken = nil
			case SERVICE_PIPEDRIVE:
				tenantSettings.PipedriveApiToken = nil
			case SERVICE_PLAID:
				tenantSettings.PlaidAccessToken = nil
			case SERVICE_PLAUSIBLE:
				tenantSettings.PlausibleApiKey = nil
				tenantSettings.PlausibleSiteId = nil
			case SERVICE_POSTHOG:
				tenantSettings.PostHogApiKey = nil
				tenantSettings.PostHogBaseUrl = nil
			case SERVICE_QUALAROO:
				tenantSettings.QualarooApiKey = nil
				tenantSettings.QualarooApiToken = nil
			case SERVICE_QUICKBOOKS:
				tenantSettings.QuickBooksClientId = nil
				tenantSettings.QuickBooksClientSecret = nil
				tenantSettings.QuickBooksRealmId = nil
				tenantSettings.QuickBooksRefreshToken = nil
			case SERVICE_RECHARGE:
				tenantSettings.RechargeApiToken = nil
			case SERVICE_RECRUITEE:
				tenantSettings.RecruiteeCompanyId = nil
				tenantSettings.RecruiteeApiKey = nil
			case SERVICE_RECURLY:
				tenantSettings.RecurlyApiKey = nil
			case SERVICE_RETENTLY:
				tenantSettings.RetentlyApiToken = nil
			case SERVICE_SALESFORCE:
				tenantSettings.SalesforceClientId = nil
				tenantSettings.SalesforceClientSecret = nil
				tenantSettings.SalesforceRefreshToken = nil
			case SERVICE_SALESLOFT:
				tenantSettings.SalesloftApiKey = nil
			case SERVICE_SENDGRID:
				tenantSettings.SendgridApiKey = nil
			case SERVICE_SENTRY:
				tenantSettings.SentryProject = nil
				tenantSettings.SentryHost = nil
				tenantSettings.SentryAuthenticationToken = nil
				tenantSettings.SentryOrganization = nil
			case SERVICE_SLACK:
				tenantSettings.SlackApiToken = nil
				tenantSettings.SlackChannelFilter = nil
				tenantSettings.SlackLookbackWindow = nil
			case SERVICE_STRIPE:
				tenantSettings.StripeAccountId = nil
				tenantSettings.StripeSecretKey = nil
			case SERVICE_SURVEYSPARROW:
				tenantSettings.SurveySparrowAccessToken = nil
			case SERVICE_SURVEYMONKEY:
				tenantSettings.SurveyMonkeyAccessToken = nil
			case SERVICE_TALKDESK:
				tenantSettings.TalkdeskApiKey = nil
			case SERVICE_TIKTOK:
				tenantSettings.TikTokAccessToken = nil
			case SERVICE_TODOIST:
				tenantSettings.TodoistApiToken = nil
			case SERVICE_TYPEFORM:
				tenantSettings.TypeformApiToken = nil
			case SERVICE_VITTALLY:
				tenantSettings.VittallyApiKey = nil
			case SERVICE_WRIKE:
				tenantSettings.WrikeAccessToken = nil
			case SERVICE_XERO:
				tenantSettings.XeroClientId = nil
				tenantSettings.XeroClientSecret = nil
				tenantSettings.XeroTenantId = nil
				tenantSettings.XeroScopes = nil
			case SERVICE_ZENDESK_SUPPORT:
				tenantSettings.ZendeskAPIKey = nil
				tenantSettings.ZendeskSubdomain = nil
				tenantSettings.ZendeskAdminEmail = nil
			case SERVICE_ZENDESK_CHAT:
				tenantSettings.ZendeskChatSubdomain = nil
				tenantSettings.ZendeskChatAccessKey = nil
			case SERVICE_ZENDESK_TALK:
				tenantSettings.ZendeskTalkSubdomain = nil
				tenantSettings.ZendeskTalkAccessKey = nil
			case SERVICE_ZENDESK_SELL:
				tenantSettings.ZendeskSellApiToken = nil
			case SERVICE_ZENDESK_SUNSHINE:
				tenantSettings.ZendeskSunshineSubdomain = nil
				tenantSettings.ZendeskSunshineApiToken = nil
				tenantSettings.ZendeskSunshineEmail = nil
			case SERVICE_ZENEFITS:
				tenantSettings.ZenefitsToken = nil
			case SERVICE_MIXPANEL:
				tenantSettings.MixpanelUsername = nil
				tenantSettings.MixpanelSecret = nil
				tenantSettings.MixpanelProjectSecret = nil
				tenantSettings.MixpanelProjectId = nil
				tenantSettings.MixpanelProjectTimezone = nil
				tenantSettings.MixpanelRegion = nil
			}
		}

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, nil, qr.Error
		}

		activeServices, err := s.GetServiceActivations(tenantName)
		if err != nil {
			return nil, nil, fmt.Errorf("ClearIntegrationData: %v", err)
		}
		return qr.Result.(*entity.TenantSettings), activeServices, nil
	}
}
