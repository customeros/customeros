package interfaces

import (
	"context"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type MeetingService interface {
	GetCalendarAvailabilityForTenant(ctx context.Context, startTime time.Time, endTime time.Time, duration int, timezone string) (*CalendarAvailability, error)
	GetCalendarAvailabilityForEmail(ctx context.Context, email string, startTime time.Time, endTime time.Time, duration int, timezone string) (*CalendarAvailability, error)

	GetUserCalendarAvailability(ctx context.Context, email string) (*postgres_entity.UserCalendarAvailability, error)
	SaveUserCalendarAvailability(ctx context.Context, availability *postgres_entity.UserCalendarAvailability) (*postgres_entity.UserCalendarAvailability, error)
	SetDefaultUserCalendarAvailability(ctx context.Context, email string) (*postgres_entity.UserCalendarAvailability, error)
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
