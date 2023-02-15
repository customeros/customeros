package dto

type TenantSettingsResponseDTO struct {
	HubspotExists    bool `json:"hubspotExists"`
	ZendeskExists    bool `json:"zendeskExists"`
	SmartSheetExists bool `json:"smartSheetExists"`
	JiraExists       bool `json:"jiraExists"`
	TrelloExists     bool `json:"trelloExists"`
}
