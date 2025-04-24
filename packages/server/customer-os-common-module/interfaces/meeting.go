package interfaces

import (
	"context"
	"time"

	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type MeetingService interface {
	GetUserCalendarAvailability(ctx context.Context, email string) (*postgresEntity.UserCalendarAvailability, error)
	SaveUserCalendarAvailability(ctx context.Context, availability *postgresEntity.UserCalendarAvailability) (*postgresEntity.UserCalendarAvailability, error)
	SetDefaultUserCalendarAvailability(ctx context.Context, email string) (*postgresEntity.UserCalendarAvailability, error)
}

// TimeSlot represents a time slot in calendar availability
type TimeSlot struct {
	StartTime   time.Time
	EndTime     time.Time
	IsAvailable bool
}

// CalendarAvailability represents the availability of users for a given time range
type CalendarAvailability struct {
	TimeSlots []TimeSlot
}
