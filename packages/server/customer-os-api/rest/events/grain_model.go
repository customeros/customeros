package events

import "time"

// RecordingResponse represents the top-level response structure
type GrainRecordingData struct {
	Data   GrainRecording `json:"data"`
	Type   string         `json:"type"`
	UserID string         `json:"user_id"`
}

// RecordingData contains the details of the recording
type GrainRecording struct {
	EndDatetime         time.Time          `json:"end_datetime"`
	ICalUID             string             `json:"ical_uid"`
	ID                  string             `json:"id"`
	IntelligenceNotesMD string             `json:"intelligence_notes_md"`
	Owners              []string           `json:"owners"`
	Participants        []GrainParticipant `json:"participants"`
	PublicThumbnailURL  string             `json:"public_thumbnail_url"`
	PublicURL           string             `json:"public_url"`
	StartDatetime       time.Time          `json:"start_datetime"`
	Tags                []string           `json:"tags"`
	Title               string             `json:"title"`
	URL                 string             `json:"url"`
}

// Participant represents an individual participant in the recording
type GrainParticipant struct {
	Name              string  `json:"name"`
	Scope             string  `json:"scope"`
	Email             *string `json:"email,omitempty"`
	ConfirmedAttendee bool    `json:"confirmed_attendee"`
}
