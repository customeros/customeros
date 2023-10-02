package entity

import (
	"github.com/google/uuid"
	"time"
)

type RawCalendarEvent struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`

	ExternalSystem string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	TenantName     string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	UsernameSource string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	CalendarId     string `gorm:"size:255;not null;index:idx_raw_email_external_system"`

	ProviderId string `gorm:"size:255;not null;index:idx_raw_email_external_system"`
	ICalUID    string `gorm:"size:255;not null;"`

	SentToEventStoreState  string  `gorm:"size:50;not null"`
	SentToEventStoreReason *string `gorm:"type:text"`
	SentToEventStoreError  *string `gorm:"type:text"`

	Data string `gorm:"type:text"`
}

func (RawCalendarEvent) TableName() string {
	return "raw_calendar_event"
}
