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

	// Calendar operations
	ListCalendars(ctx context.Context, email string) ([]*NylasCalendar, error)
	GetDefaultCalendar(ctx context.Context, email string) (*NylasCalendar, error)
	GetCalendarAvailability(ctx context.Context, email, calendarID string, startTime, endTime time.Time, bufferBefore, bufferAfter int) (*NylasAvailabilityResponse, error)

	// Meeting operations
	CreateEvent(ctx context.Context, meetingData NylasCreateEventRequest, hostEmail, calendarID string, notifyParticipants bool) (*NylasEvent, error)
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

type NylasAvailabilityParticipant struct {
	Email       string   `json:"email"`
	CalendarIds []string `json:"calendar_ids,omitempty"`
	OpenHours   []struct {
		Days     []int    `json:"days"`
		Timezone string   `json:"timezone"`
		Start    string   `json:"start"`
		End      string   `json:"end"`
		Exdates  []string `json:"exdates"`
	} `json:"open_hours,omitempty"`
}

type NylasAvailabilityRules struct {
	AvailabilityMethod string `json:"availability_method"`
	Buffer             struct {
		Before int `json:"before"`
		After  int `json:"after"`
	} `json:"buffer"`
	TentativeAsBusy  bool `json:"tentative_as_busy"`
	DefaultOpenHours []struct {
		Days     []int    `json:"days"`
		Timezone string   `json:"timezone"`
		Start    string   `json:"start"`
		End      string   `json:"end"`
		Exdates  []string `json:"exdates"`
	} `json:"default_open_hours"`
}

type NylasAvailabilityRequest struct {
	Participants      []NylasAvailabilityParticipant `json:"participants"`
	StartTime         int64                          `json:"start_time"`
	EndTime           int64                          `json:"end_time"`
	IntervalMinutes   int                            `json:"interval_minutes"`
	DurationMinutes   int                            `json:"duration_minutes"`
	RoundTo           int                            `json:"round_to"`
	AvailabilityRules NylasAvailabilityRules         `json:"availability_rules"`
}

type NylasAvailabilityResponse struct {
	RequestID string `json:"request_id"`
	Data      struct {
		Order     []string `json:"order"`
		TimeSlots []struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
		} `json:"time_slots"`
	} `json:"data"`
}

type NylasCreateEventRequest struct {
	Busy        bool   `json:"busy"`
	Title       string `json:"title"`
	Description string `json:"description"`
	When        struct {
		StartTime     int64  `json:"start_time"`
		EndTime       int64  `json:"end_time"`
		StartTimezone string `json:"start_timezone"`
		EndTimezone   string `json:"end_timezone"`
	} `json:"when"`
	Location     string                  `json:"location,omitempty"`
	Participants []NylasEventParticipant `json:"participants"`
	Resources    []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"resources,omitempty"`
	Recurrence []string `json:"recurrence,omitempty"`
}

type NylasEventParticipant struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type NylasEvent struct {
	RequestID string `json:"request_id"`
	Data      struct {
		Busy       bool   `json:"busy"`
		CalendarID string `json:"calendar_id"`
		HtmlLink   string `json:"html_link"`
	} `json:"data"`
}
