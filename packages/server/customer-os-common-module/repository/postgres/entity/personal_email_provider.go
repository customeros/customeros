package entity

import "github.com/google/uuid"

type PersonalEmailProvider struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProviderName   string    `gorm:"size:255;not null;"`
	ProviderDomain string    `gorm:"size:255;not null;index:idx_provider_domain"`
}

func (PersonalEmailProvider) TableName() string {
	return "personal_email_provider"
}
