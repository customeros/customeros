package flows

import "time"

// RecordingRes the top-level response structure
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
	OwnersStr           string             `json:"owners"`
	Owners              []string           `json:"ownersArr"`
	ParticipantsStr     string             `json:"participants"`
	Participants        []GrainParticipant `json:"participantsArray"`
	PublicThumbnailURL  string             `json:"public_thumbnail_url"`
	PublicURL           string             `json:"public_url"`
	StartDatetime       time.Time          `json:"start_datetime"`
	TagsStr             string             `json:"tags"`
	Tags                []string           `json:"tagsArray"`
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
