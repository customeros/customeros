package interfaces

import (
	"context"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type NylasProvider string

const (
	NylasProviderGoogle NylasProvider = "google"
)

type NylasService interface {
	// Access operations
	GrantAccess(ctx context.Context, email, userId, refreshToken string, nylasProvider NylasProvider) (*postgres_entity.NylasGrant, error)
	RevokeAccess(ctx context.Context, email string) error
	GetGrant(ctx context.Context, email string) (*postgres_entity.NylasGrant, error)

	// Calendar management
	ListCalendars(ctx context.Context, email string) ([]*Calendar, error)
	GetDefaultCalendar(ctx context.Context, email string) (*Calendar, error)
	GetCalendar(ctx context.Context, calendarID string) (*Calendar, error)

	// Calendar operations
	CreateEvent(ctx context.Context, calendarID string, event *CalendarEvent) (*CalendarEvent, error)
	UpdateEvent(ctx context.Context, calendarID string, eventID string, event *CalendarEvent) (*CalendarEvent, error)
	DeleteEvent(ctx context.Context, calendarID string, eventID string) error
	GetEvent(ctx context.Context, calendarID string, eventID string) (*CalendarEvent, error)
	ListEvents(ctx context.Context, calendarID string, startTime, endTime time.Time) ([]*CalendarEvent, error)
}

type CalendarEvent struct {
	ID           string
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	Location     string
	Participants []Participant
	Recurrence   *Recurrence
	Status       string
}

type Participant struct {
	Email  string
	Name   string
	Status string // accepted, declined, tentative
}

type Recurrence struct {
	Frequency string // daily, weekly, monthly, yearly
	Interval  int
	Until     time.Time
}

type Calendar struct {
	ID          string
	Name        string
	Description string
	IsPrimary   bool
}
