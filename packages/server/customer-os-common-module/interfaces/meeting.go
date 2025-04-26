package interfaces

import (
	"context"
	"time"

	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

// DaySlot represents a day's availability
type DaySlot struct {
	Date        time.Time
	TimeSlots   []TimeSlot
	IsAvailable bool
}

// TimeSlot represents a time slot in calendar availability
type TimeSlot struct {
	StartTime   time.Time
	EndTime     time.Time
	IsAvailable bool
}

// CalendarAvailabilityResult represents the result of calendar availability calculation
type CalendarAvailabilityResult struct {
	Days []*DaySlot
}

type MeetingService interface {
	GetUserCalendarAvailability(ctx context.Context, email string) (*postgresEntity.UserCalendarAvailability, error)
	SaveUserCalendarAvailability(ctx context.Context, availability *postgresEntity.UserCalendarAvailability) (*postgresEntity.UserCalendarAvailability, error)
	SetDefaultUserCalendarAvailability(ctx context.Context, email, timezone string) (*postgresEntity.UserCalendarAvailability, error)

	// GetCalendarAvailability returns the available time slots for a meeting booking event
	// It takes into account:
	// - The calendars of all participants
	// - Working hours of each participant
	// - Meeting booking event rules (buffer between meetings, min notice time, etc.)
	// - Timezone conversion
	GetCalendarAvailability(ctx context.Context, meetingBookingEventID string, startTime time.Time, endTime time.Time, timezone string) (*CalendarAvailabilityResult, error)
}
