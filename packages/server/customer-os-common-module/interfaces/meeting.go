package interfaces

import (
	"context"
	"time"
)

type MeetingService interface {
	GetCalendarAvailabilityForTenant(ctx context.Context, startTime time.Time, endTime time.Time, duration int, timezone string) (*CalendarAvailability, error)
	GetCalendarAvailabilityForEmail(ctx context.Context, email string, startTime time.Time, endTime time.Time, duration int, timezone string) (*CalendarAvailability, error)
}

// TimeSlot represents a time slot in calendar availability
type TimeSlot struct {
	StartTime   time.Time
	EndTime     time.Time
	IsAvailable bool
}

// CalendarAvailability represents the availability of users for a given time range
type CalendarAvailability struct {
	TimeSlots      []TimeSlot
	TotalUsers     int
	AvailableUsers int
}
