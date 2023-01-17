package dto

type TenantSettingsResponseDTO struct {
	HubspotExists    bool `json:"hubspotExists"`
	ZendeskExists    bool `json:"zendeskExists"`
	SmartSheetExists bool `json:"smartSheetExists"`
}
