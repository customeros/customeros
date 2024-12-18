package model

type OAuthUserSettingsResponse struct {
	Provider           string `json:"provider"`
	Email              string `json:"email"`
	NeedsManualRefresh bool   `json:"needsManualRefresh"`
	Type               string `json:"type"`
}
