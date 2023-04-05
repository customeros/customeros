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
		}

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}
