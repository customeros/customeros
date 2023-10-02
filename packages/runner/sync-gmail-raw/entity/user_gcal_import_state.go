package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserGCalImportState struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantName string    `gorm:"size:255;not null"`
	Username   string    `gorm:"size:255;not null"`
	CalendarId string    `gorm:"size:255;not null"`

	SyncToken  string    `gorm:"size:255;not null"`
	PageToken  string    `gorm:"type:text;not null"`
	MaxResults int64     `gorm:"size:255;not null"`
	TimeMin    time.Time `gorm:"size:255;not null"`
	TimeMax    time.Time `gorm:"size:255;not null"`
}

func (UserGCalImportState) TableName() string {
	return "user_gcal_import_state"
}
