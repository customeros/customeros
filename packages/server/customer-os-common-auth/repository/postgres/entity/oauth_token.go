package entity

import (
	"time"
)

const (
	ProviderGoogle = "google"
)

type OAuthTokenEntity struct {
	Provider                  string    `gorm:"primaryKey;autoIncrement:false;column:provider;size:255;not null"`
	PlayerIdentityId          string    `gorm:"primaryKey;autoIncrement:false;column:player_identity_id;size:255;not null"`
	TenantName                string    `gorm:"index:,size:255;not null"`
	EmailAddress              string    `gorm:"column:email_address;size:255;"`
	AccessToken               string    `gorm:"column:access_token;type:text"`
	RefreshToken              string    `gorm:"column:refresh_token;type:text"`
	NeedsManualRefresh        bool      `gorm:"column:needs_manual_refresh;default:false;"`
	IdToken                   string    `gorm:"column:id_token;type:text"`
	ExpiresAt                 time.Time `gorm:"column:expires_at;type:timestamp;"`
	Scope                     string    `gorm:"column:scope;type:text"`
	GmailSyncEnabled          bool      `gorm:"column:gmail_sync_enabled;default:false;"`
	GoogleCalendarSyncEnabled bool      `gorm:"column:google_calendar_sync_enabled;default:false;"`
}

func (OAuthTokenEntity) TableName() string {
	return "oauth_token"
}
