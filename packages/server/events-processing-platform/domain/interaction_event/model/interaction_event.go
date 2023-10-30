package model

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type InteractionEvent struct {
	Tenant             string                       `json:"tenant"`
	ID                 string                       `json:"id"`
	Source             commonmodel.Source           `json:"source"`
	CreatedAt          time.Time                    `json:"createdAt"`
	UpdatedAt          time.Time                    `json:"updatedAt"`
	Content            string                       `json:"content"`
	ContentType        string                       `json:"contentType"`
	Channel            string                       `json:"channel"`
	ChannelData        string                       `json:"channelData"`
	EventType          string                       `json:"eventType"`
	Identifier         string                       `json:"identifier"`
	Summary            string                       `json:"summary"`
	ActionItems        []string                     `json:"actionItems"`
	ExternalSystems    []commonmodel.ExternalSystem `json:"externalSystem"`
	BelongsToIssueId   string                       `json:"belongsToIssueId,omitempty"`
	BelongsToSessionId string                       `json:"belongsToSessionId,omitempty"`
	Hide               bool                         `json:"hide"`
	Sender             Sender                       `json:"sender"`
	Receivers          []Receiver                   `json:"receivers"`
}

type Sender struct {
	Participant  commonmodel.Participant `json:"participant"`
	RelationType string                  `json:"relationType"`
}

type Receiver struct {
	Participant  commonmodel.Participant `json:"participant"`
	RelationType string                  `json:"relationType"`
}

func (s Sender) Available() bool {
	return s.Participant.ID != "" && s.Participant.ParticipantType != ""
}

func (r Receiver) Available() bool {
	return r.Participant.ID != "" && r.Participant.ParticipantType != ""
}
