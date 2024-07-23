package entity

import (
	"github.com/google/uuid"
)

type TenantSettingsEmailExclusion struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant         string    `gorm:"size:255;not null;"`
	ExcludeSubject *string   `gorm:"size:255;"`
	ExcludeBody    *string   `gorm:"size:255;"`
}

func (TenantSettingsEmailExclusion) TableName() string {
	return "tenant_settings_email_exclusion"
}
