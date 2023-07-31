package entity

import "github.com/google/uuid"

type RawEmail struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	ExternalSystem string `gorm:"size:255;not null"`
	TenantName     string `gorm:"size:255;not null"`
	UsernameSource string `gorm:"size:255;not null"`
	MessageId      string `gorm:"size:255;not null"`

	SentToEventStore bool `gorm:"size:255;not null"`

	Data string `gorm:"type:text"`
}

func (RawEmail) TableName() string {
	return "raw_email"
}
