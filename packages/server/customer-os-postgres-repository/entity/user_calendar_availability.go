package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

// DayAvailability represents the availability for a single day
type DayAvailability struct {
	Enabled   bool   `json:"enabled" gorm:"type:boolean;default:false"`
	StartHour string `json:"startHour" gorm:"type:varchar(5)"` // Format: "HH:MM"
	EndHour   string `json:"endHour" gorm:"type:varchar(5)"`   // Format: "HH:MM"
}

// UserCalendarAvailability represents a user's available hours for calendar bookings
type UserCalendarAvailability struct {
	ID        string    `gorm:"type:varchar(21);primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp"`
	Tenant    string    `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	Email     string    `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	Timezone  string    `gorm:"size:50"` // e.g., "America/New_York"

	// One field for each day of the week
	Monday    DayAvailability `gorm:"type:jsonb"`
	Tuesday   DayAvailability `gorm:"type:jsonb"`
	Wednesday DayAvailability `gorm:"type:jsonb"`
	Thursday  DayAvailability `gorm:"type:jsonb"`
	Friday    DayAvailability `gorm:"type:jsonb"`
	Saturday  DayAvailability `gorm:"type:jsonb"`
	Sunday    DayAvailability `gorm:"type:jsonb"`
}

func (UserCalendarAvailability) TableName() string {
	return "user_calendar_availability"
}

// BeforeCreate hook to ensure ID has the correct prefix
func (u *UserCalendarAvailability) BeforeCreate(tx *gorm.DB) error {
	u.ID = utils.GenerateNanoIdWithPrefix("uca", 16)
	return nil
}
