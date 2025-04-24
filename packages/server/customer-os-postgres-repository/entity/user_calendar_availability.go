package postgres_entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

// DayAvailability represents the availability for a single day
type DayAvailability struct {
	Enabled   bool   `json:"enabled" gorm:"type:boolean;default:false"`
	StartHour string `json:"startHour" gorm:"type:varchar(5)"` // Format: "HH:MM"
	EndHour   string `json:"endHour" gorm:"type:varchar(5)"`   // Format: "HH:MM"
}

// standardizeTimeFormat converts time strings like "8:10" to "08:10"
func standardizeTimeFormat(timeStr string) (string, error) {
	if timeStr == "" {
		return "", nil
	}

	// Special case for end of day
	if timeStr == "24:00" {
		return "24:00", nil
	}

	// Try parsing with different formats
	formats := []string{"15:04", "3:04PM", "3:04 PM", "3PM", "15", "3"}
	var t time.Time
	var err error

	for _, format := range formats {
		t, err = time.Parse(format, timeStr)
		if err == nil {
			return t.Format("15:04"), nil
		}
	}

	// If single digit hour provided, try prefixing with 0
	if !strings.Contains(timeStr, ":") {
		timeStr = timeStr + ":00"
	}
	if len(timeStr) == 4 && timeStr[1] == ':' {
		timeStr = "0" + timeStr
	}

	// Final attempt with standardized format
	t, err = time.Parse("15:04", timeStr)
	if err != nil {
		return "", fmt.Errorf("invalid time format: %s", timeStr)
	}

	return t.Format("15:04"), nil
}

// compareTimeStrings compares two time strings in "HH:MM" format
// returns true if start is before end
func compareTimeStrings(start, end string) bool {
	// Parse times
	startParts := strings.Split(start, ":")
	endParts := strings.Split(end, ":")

	startHour, _ := strconv.Atoi(startParts[0])
	startMin, _ := strconv.Atoi(startParts[1])
	endHour, _ := strconv.Atoi(endParts[0])
	endMin, _ := strconv.Atoi(endParts[1])

	// Compare hours first
	if startHour < endHour {
		return true
	}
	if startHour > endHour {
		return false
	}
	// If hours are equal, compare minutes
	return startMin < endMin
}

// Validate checks if the day availability is valid
func (d *DayAvailability) Validate() error {
	if !d.Enabled {
		return nil
	}

	// Validate start time
	start, err := standardizeTimeFormat(d.StartHour)
	if err != nil {
		return err
	}
	d.StartHour = start

	// Validate end time
	end, err := standardizeTimeFormat(d.EndHour)
	if err != nil {
		return err
	}
	d.EndHour = end

	// Skip time comparison if either time is empty
	if start == "" || end == "" {
		return nil
	}

	// Validate start time is before end time
	if !compareTimeStrings(start, end) {
		return fmt.Errorf("start time %s must be before end time %s", start, end)
	}

	return nil
}

// Scan implements the sql.Scanner interface for DayAvailability
func (d *DayAvailability) Scan(value interface{}) error {
	if value == nil {
		*d = DayAvailability{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, d)
}

// Value implements the driver.Valuer interface for DayAvailability
func (d DayAvailability) Value() (driver.Value, error) {
	if err := d.standardizeAndValidate(); err != nil {
		return nil, err
	}
	return json.Marshal(d)
}

// UserCalendarAvailability represents a user's available hours for calendar bookings
type UserCalendarAvailability struct {
	ID        string    `gorm:"type:varchar(21);primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp"`
	Tenant    string    `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	Email     string    `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	Timezone  string    `gorm:"size:50;not null"` // e.g., "America/New_York"

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

// Validate performs validation on the UserCalendarAvailability entity
func (u *UserCalendarAvailability) Validate() error {
	// Validate timezone
	if u.Timezone == "" {
		u.Timezone = "Etc/UTC"
	}
	// Trim any whitespace that might have been added
	u.Timezone = strings.TrimSpace(u.Timezone)
	_, err := time.LoadLocation(u.Timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %s, %s", u.Timezone, err.Error())
	}

	// Validate each day's availability
	days := []struct {
		name string
		day  *DayAvailability
	}{
		{"Monday", &u.Monday},
		{"Tuesday", &u.Tuesday},
		{"Wednesday", &u.Wednesday},
		{"Thursday", &u.Thursday},
		{"Friday", &u.Friday},
		{"Saturday", &u.Saturday},
		{"Sunday", &u.Sunday},
	}

	for _, d := range days {
		if err := d.day.standardizeAndValidate(); err != nil {
			return fmt.Errorf("%s: %v", d.name, err)
		}
	}

	return nil
}

// BeforeCreate hook to ensure ID has the correct prefix and validate
func (u *UserCalendarAvailability) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		id := utils.GenerateNanoIdWithPrefix("uca", 16)
		u.ID = id
	} else if !strings.HasPrefix(u.ID, "uca_") {
		u.ID = "uca_" + u.ID
	}
	return u.Validate()
}

// BeforeUpdate hook to update the UpdatedAt timestamp and validate
func (u *UserCalendarAvailability) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return u.Validate()
}

// standardizeAndValidate checks if the day availability is valid and standardizes time formats
func (d *DayAvailability) standardizeAndValidate() error {
	if !d.Enabled {
		d.StartHour = ""
		d.EndHour = ""
		return nil
	}

	// Validate start time
	start, err := standardizeTimeFormat(d.StartHour)
	if err != nil {
		return err
	}
	d.StartHour = start

	// Validate end time
	end, err := standardizeTimeFormat(d.EndHour)
	if err != nil {
		return err
	}
	d.EndHour = end

	// Skip time comparison if either time is empty
	if start == "" || end == "" {
		return nil
	}

	// Validate start time is before end time
	if !compareTimeStrings(start, end) {
		return fmt.Errorf("start time %s must be before end time %s", start, end)
	}

	return nil
}
