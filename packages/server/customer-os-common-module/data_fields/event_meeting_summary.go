package data_fields

import (
	"encoding/json"
	"time"
)

type MeetingSummaryEvent struct {
	MeetingID         string     `json:"meetingId"`
	ParticipantEmails *[]string  `json:"participantEmails,omitempty"`
	Content           *string    `json:"content,omitempty"`
	Timestamp         *time.Time `json:"timestamp,omitempty"`
}

func (m MeetingSummaryEvent) Type() string {
	return "MeetingSummaryEvent"
}

func (m MeetingSummaryEvent) ToString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MeetingSummaryEventFromString(s string) (*MeetingSummaryEvent, error) {
	var event MeetingSummaryEvent
	if err := json.Unmarshal([]byte(s), &event); err != nil {
		return nil, err
	}
	return &event, nil
}
