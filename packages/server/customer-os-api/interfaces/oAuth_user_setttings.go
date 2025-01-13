package cosapi_interfaces

import "context"

type OAuthUserSettingsService interface {
	GetTenantOAuthUserSettings(ctx context.Context, tenant string) ([]*OAuthUserSettingsResponse, error)
}

type OAuthUserSettingsResponse struct {
	Provider           string `json:"provider"`
	Email              string `json:"email"`
	NeedsManualRefresh bool   `json:"needsManualRefresh"`
	Type               string `json:"type"`
}
