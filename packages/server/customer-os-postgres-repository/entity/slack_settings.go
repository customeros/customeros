package entity

import (
	"github.com/google/uuid"
)

type SlackSettingsEntity struct {
	Id           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantName   string    `gorm:"index:idx_tenant_uk;size:255;not null"`
	AppId        string    `gorm:"size:255;"`
	AuthedUserId string    `gorm:"size:255;"`
	Scope        string    `gorm:"size:255;"`
	TokenType    string    `gorm:"size:255;"`
	AccessToken  string    `gorm:"size:255;"`
	BotUserId    string    `gorm:"size:255;"`
	TeamId       string    `gorm:"size:255;"`
}

func (SlackSettingsEntity) TableName() string {
	return "slack_settings"
}
