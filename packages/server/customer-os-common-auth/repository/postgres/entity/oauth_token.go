package entity

import (
	"github.com/google/uuid"
	"time"
)

const (
	ProviderGoogle = "google"
)

type OAuthTokenEntity struct {
	ID               uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Provider         string    `gorm:"index:,unique,composite:uniquePlayerProviderToken;column:provider;size:255;not null"`
	TenantName       string    `gorm:"index:,unique,composite:uniquePlayerProviderToken;column:tenant_name;size:255;not null"`
	PlayerIdentityId string    `gorm:"index:,unique,composite:uniquePlayerProviderToken;column:player_identity_id;size:255;not null"`
	EmailAddress     string    `gorm:"column:email_address;size:255;"`
	AccessToken      string    `gorm:"column:access_token;size:255;"`
	RefreshToken     string    `gorm:"column:refresh_token;size:255;"`
	ExpiresAt        time.Time `gorm:"column:expires_at;type:timestamp;"`
	Scope            string    `gorm:"column:scope;size:255;"`
	EnabledForSync   bool      `gorm:"column:enabled_for_sync;default:false;"`
}

func (OAuthTokenEntity) TableName() string {
	return "oauth_token"
}
