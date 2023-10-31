package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserGmailImportStateHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EntityId  uuid.UUID `gorm:"not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`

	Tenant    string           `gorm:"size:255;not null"`
	Username  string           `gorm:"size:255;not null"`
	State     GmailImportState `gorm:"size:50;not null"`
	StartDate *time.Time       `gorm:""`
	StopDate  *time.Time       `gorm:""`
	Active    bool             `gorm:"not null"`
	Cursor    string           `gorm:"size:255;not null"`
}

func (UserGmailImportStateHistory) TableName() string {
	return "user_gmail_import_state_history"
}
