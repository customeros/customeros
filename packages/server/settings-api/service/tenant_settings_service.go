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

		}

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}
