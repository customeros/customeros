package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserWorkingSchedule struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	Tenant    string    `gorm:"size:255;not null;"`

	UserId    string `gorm:"type:varchar(255);not null"`
	DayRange  string `gorm:"type:varchar(7);not null"` // E.g., "Mon-Fri" or "Thu-Fri"
	StartHour string `gorm:"type:varchar(5);not null"` // E.g., 09:00
	EndHour   string `gorm:"type:varchar(5);not null"` // E.g., 18:00
}

func (UserWorkingSchedule) TableName() string {
	return "user_working_schedule"
}
