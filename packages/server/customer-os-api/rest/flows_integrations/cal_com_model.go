package integrations

import "time"

// WebhookPayload represents the root structure of all Cal.com webhooks
type CalDotComPayload struct {
	TriggerEvent string    `json:"triggerEvent"`
	CreatedAt    time.Time `json:"createdAt"`
	Payload      Booking   `json:"payload"`
}

// Booking represents the main booking data structure
type Booking struct {
	BookerURL           string              `json:"bookerUrl"`
	Type                string              `json:"type"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	AdditionalNotes     string              `json:"additionalNotes,omitempty"`
	CustomInputs        map[string]string   `json:"customInputs"`
	StartTime           time.Time           `json:"startTime"`
	EndTime             time.Time           `json:"endTime"`
	Organizer           Organizer           `json:"organizer"`
	Responses           BookingResponses    `json:"responses"`
	UserFieldsResponses UserFieldsResponses `json:"userFieldsResponses"`
	Attendees           []Attendee          `json:"attendees"`
	Location            string              `json:"location"`
	DestinationCalendar []Calendar          `json:"destinationCalendar"`
	EventTypeID         int                 `json:"eventTypeId"`
	UID                 string              `json:"uid"`
	BookingID           int                 `json:"bookingId"`
	EventTitle          string              `json:"eventTitle"`
	EventDescription    string              `json:"eventDescription"`
	Price               int                 `json:"price"`
	Currency            string              `json:"currency"`
	Status              string              `json:"status"`
	Length              int                 `json:"length"`
	ICalUID             string              `json:"iCalUID"`
	ICalSequence        int                 `json:"iCalSequence"`

	// Fields specific to cancellation
	CancellationReason string `json:"cancellationReason,omitempty"`
	CancelledBy        string `json:"cancelledBy,omitempty"`

	// Fields specific to rescheduling
	RescheduleID        int       `json:"rescheduleId,omitempty"`
	RescheduleUID       string    `json:"rescheduleUid,omitempty"`
	RescheduleStartTime time.Time `json:"rescheduleStartTime,omitempty"`
	RescheduleEndTime   time.Time `json:"rescheduleEndTime,omitempty"`
	RescheduledBy       string    `json:"rescheduledBy,omitempty"`
}

type Organizer struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	TimeZone   string   `json:"timeZone"`
	Language   Language `json:"language"`
	TimeFormat string   `json:"timeFormat"`
	UTCOffset  int      `json:"utcOffset"`
}

type Attendee struct {
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	TimeZone  string   `json:"timeZone"`
	Language  Language `json:"language"`
	UTCOffset int      `json:"utcOffset"`
}

type Language struct {
	Locale string `json:"locale"`
}

type Calendar struct {
	ID           int    `json:"id"`
	Integration  string `json:"integration"`
	ExternalID   string `json:"externalId"`
	PrimaryEmail string `json:"primaryEmail"`
	UserID       int    `json:"userId"`
	EventTypeID  *int   `json:"eventTypeId"`
	CredentialID int    `json:"credentialId"`
}

type BookingResponses struct {
	Name             ResponseField `json:"name"`
	Email            ResponseField `json:"email"`
	Notes            ResponseField `json:"notes,omitempty"`
	Guests           GuestsField   `json:"guests"`
	Location         LocationField `json:"location"`
	PhoneNumber      ResponseField `json:"phone-number"`
	RescheduleReason ResponseField `json:"rescheduleReason,omitempty"`
}

type ResponseField struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	IsHidden bool   `json:"isHidden,omitempty"`
}

type LocationField struct {
	Label    string        `json:"label"`
	Value    LocationValue `json:"value"`
	IsHidden bool          `json:"isHidden,omitempty"`
}

type LocationValue struct {
	Value       string `json:"value"`
	OptionValue string `json:"optionValue"`
}

type GuestsField struct {
	Label    string   `json:"label"`
	Value    []string `json:"value"`
	IsHidden bool     `json:"isHidden,omitempty"`
}

type UserFieldsResponses struct {
	PhoneNumber ResponseField `json:"phone-number"`
}
