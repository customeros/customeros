package model

type OAuthUserSettingsResponse struct {
	TenantName             string `json:"tenantName"`
	EmailAddress           string `json:"emailAddress"`
	GoogleOAuthSyncEnabled bool   `json:"GoogleOAuthSyncEnabled"`
}
