package dto

type TenantSettingsDTO struct {
	Id                   string `json:"id"`
	HubspotPrivateAppKey string `json:"hubspotPrivateAppKey"`
	ZendeskAPIKey        string `json:"zendeskAPIKey"`
	ZendeskSubdomain     string `json:"zendeskSubdomain"`
	ZendeskAdminEmail    string `json:"zendeskAdminEmail"`
}
