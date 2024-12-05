package data_fields

import (
	"time"
)

type MeetingSummaryEvent struct {
	ParticipantEmails *[]string  `json:"participantEmails"`
	Content           *string    `json:"content"`
	Timestamp         *time.Time `json:"timestamp"`
}

func (m *MeetingSummaryEvent) Type() string {
	return "MeetingSummaryEvent"
}
