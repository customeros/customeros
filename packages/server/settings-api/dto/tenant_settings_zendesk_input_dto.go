package dto

type TenantSettingsZendeskDTO struct {
	ZendeskAPIKey     *string `json:"zendeskAPIKey"`
	ZendeskSubdomain  *string `json:"zendeskSubdomain"`
	ZendeskAdminEmail *string `json:"zendeskAdminEmail"`
}
