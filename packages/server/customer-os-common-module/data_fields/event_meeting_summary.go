package data_fields

import (
	"time"
)

type MeetingSummaryEvent struct {
	ParticipantEmails *[]string  `json:"participantEmails,omitempty"`
	Content           *string    `json:"content,omitempty"`
	Timestamp         *time.Time `json:"timestamp,omitempty"`
}

func (m MeetingSummaryEvent) Type() string {
	return "MeetingSummaryEvent"
}
