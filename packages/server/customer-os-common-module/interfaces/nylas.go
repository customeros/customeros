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
	ListCalendars(ctx context.Context, email string) ([]*NylasCalendar, error)
	GetDefaultCalendar(ctx context.Context, email string) (*NylasCalendar, error)

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

type NylasCalendarsResponse struct {
	RequestID  string           `json:"request_id"`
	Data       []*NylasCalendar `json:"data"`
	NextCursor string           `json:"next_cursor"`
}

type NylasCalendar struct {
	Description        string                 `json:"description"`
	HexColor           string                 `json:"hex_color"`
	HexForegroundColor string                 `json:"hex_foreground_color"`
	ID                 string                 `json:"id"`
	IsOwnedByUser      bool                   `json:"is_owned_by_user"`
	IsPrimary          bool                   `json:"is_primary"`
	Location           string                 `json:"location"`
	Metadata           map[string]interface{} `json:"metadata"`
	Name               string                 `json:"name"`
	Object             string                 `json:"object"`
	ReadOnly           bool                   `json:"read_only"`
	Timezone           string                 `json:"timezone"`
}
