package entity

import (
	"github.com/google/uuid"
)

type EmailExclusion struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant         string    `gorm:"size:255;not null;"`
	ExcludeSubject string    `gorm:"size:255;not null;index:idx_provider_domain"`
}

func (EmailExclusion) TableName() string {
	return "email_exclusion"
}
