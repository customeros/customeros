package dto

import (
	"encoding/json"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = NewMeetingRecording{}

type NewMeetingRecording struct {
	ParticipantEmails *[]string   `json:"meetingParticipantEmails,omitempty"`
	Content           *string     `json:"meetingContent,omitempty"`
	Timestamp         *time.Time  `json:"meetingTimestamp,omitempty"`
	Source            enum.Source `json:"meetingSource,omitempty"`
}

func (m NewMeetingRecording) Name() enum.AgentListenerEvent {
	return enum.EventNewMeetingRecording
}

func (m NewMeetingRecording) ToString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MeetingSummaryEventFromString(s string) (*NewMeetingRecording, error) {
	var event NewMeetingRecording
	if err := json.Unmarshal([]byte(s), &event); err != nil {
		return nil, err
	}
	return &event, nil
}
