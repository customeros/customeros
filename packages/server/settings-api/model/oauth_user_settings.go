package model

type OAuthUserSettingsResponse struct {
	GmailSyncEnabled          bool `json:"gmailSyncEnabled"`
	GoogleCalendarSyncEnabled bool `json:"googleCalendarSyncEnabled"`
}
