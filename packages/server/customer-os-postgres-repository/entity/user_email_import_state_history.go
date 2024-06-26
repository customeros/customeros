package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserEmailImportStateHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EntityId  uuid.UUID `gorm:"not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`

	Tenant    string           `gorm:"size:255;not null"`
	Username  string           `gorm:"size:255;not null"`
	Provider  string           `gorm:"size:255;not null"`
	State     EmailImportState `gorm:"size:50;not null"`
	StartDate *time.Time       `gorm:""`
	StopDate  *time.Time       `gorm:""`
	Active    bool             `gorm:"not null"`
	Cursor    string           `gorm:"size:255;not null"`
}

func (UserEmailImportStateHistory) TableName() string {
	return "user_email_import_state_history"
}
