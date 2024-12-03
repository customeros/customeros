package data_fields

import (
	"time"
)

type MeetingSummaryEvent struct {
	ParticipantEmails *[]string
	Content           *string
	Timestamp         *time.Time
}

func (m *MeetingSummaryEvent) Type() string {
	return "MeetingSummaryEvent"
}
