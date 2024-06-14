package entity

import (
	"time"
)

const (
	ProviderGoogle = "google"
)

type OAuthTokenEntity struct {
	Provider                  string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:provider;size:255;not null"`
	TenantName                string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:tenant_name;size:255;not null"`
	EmailAddress              string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:email_address;size:255;not null"`
	PlayerIdentityId          string    `gorm:"column:player_identity_id;size:255;not null"`
	AccessToken               string    `gorm:"column:access_token;type:text"`
	RefreshToken              string    `gorm:"column:refresh_token;type:text"`
	NeedsManualRefresh        bool      `gorm:"column:needs_manual_refresh;default:false;"`
	IdToken                   string    `gorm:"column:id_token;type:text"`
	ExpiresAt                 time.Time `gorm:"column:expires_at;type:timestamp;"`
	Scope                     string    `gorm:"column:scope;type:text"`
	GmailSyncEnabled          bool      `gorm:"column:gmail_sync_enabled;default:false;"`
	GoogleCalendarSyncEnabled bool      `gorm:"column:google_calendar_sync_enabled;default:false;"`
	UserId                    string    `gorm:"column:user_id;size:255;not null"`
}

func (OAuthTokenEntity) TableName() string {
	return "oauth_token"
}
