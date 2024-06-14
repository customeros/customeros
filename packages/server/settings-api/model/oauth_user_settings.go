package model

type OAuthUserSettingsResponse struct {
	Email              string `json:"email"`
	UserId             string `json:"userId"`
	NeedsManualRefresh bool   `json:"needsManualRefresh"`
	Type               string `json:"type"`
}
