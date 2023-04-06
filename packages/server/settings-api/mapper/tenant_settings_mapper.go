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

	return &responseMap
}
