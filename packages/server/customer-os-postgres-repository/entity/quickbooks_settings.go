package postgres_entity

import (
	"github.com/google/uuid"
	"time"
)

type QuickbooksSettingsEntity struct {
	Id                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant                string    `gorm:"index:idx_tenant_uk;size:255;not null"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp"`
	RealmId               string    `gorm:"column:realm_id;"`
	AccessToken           string    `gorm:"column:access_token;type:text;"`
	AccessTokenExpiresIn  int       `gorm:"column:access_expires_in;"`
	AccessTokenExpiresAt  time.Time `gorm:"column:access_token_expires_at;type:timestamp"`
	RefreshToken          string    `gorm:"column:refresh_token;type:text;"`
	RefreshTokenExpiresIn int       `gorm:"column:refresh_token_expires_in;"`
	RefreshTokenExpiresAt time.Time `gorm:"column:refresh_token_expires_at;type:timestamp"`
	RefreshTokenExpired   bool      `gorm:"column:refresh_token_expired;"`
}

func (QuickbooksSettingsEntity) TableName() string {
	return "quickbooks_settings"
}
