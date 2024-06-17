package model

type OAuthUserSettingsResponse struct {
	Provider           string `json:"provider"`
	Email              string `json:"email"`
	UserId             string `json:"userId"`
	NeedsManualRefresh bool   `json:"needsManualRefresh"`
	Type               string `json:"type"`
}
