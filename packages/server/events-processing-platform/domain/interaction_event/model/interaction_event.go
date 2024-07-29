package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"
)

type InteractionEvent struct {
	Tenant             string                  `json:"tenant"`
	ID                 string                  `json:"id"`
	Source             common.Source           `json:"source"`
	CreatedAt          time.Time               `json:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt"`
	Content            string                  `json:"content"`
	ContentType        string                  `json:"contentType"`
	Channel            string                  `json:"channel"`
	ChannelData        string                  `json:"channelData"`
	EventType          string                  `json:"eventType"`
	Identifier         string                  `json:"identifier"`
	Summary            string                  `json:"summary"`
	ActionItems        []string                `json:"actionItems"`
	ExternalSystems    []common.ExternalSystem `json:"externalSystem"`
	BelongsToIssueId   string                  `json:"belongsToIssueId,omitempty"`
	BelongsToSessionId string                  `json:"belongsToSessionId,omitempty"`
	Hide               bool                    `json:"hide"`
	Sender             Sender                  `json:"sender"`
	Receivers          []Receiver              `json:"receivers"`
}

func (e InteractionEvent) SameData(fields InteractionEventDataFields, externalSystem common.ExternalSystem) bool {
	if !externalSystem.Available() {
		return false
	}

	if externalSystem.Available() && !e.HasExternalSystem(externalSystem) {
		return false
	}

	if e.Source.SourceOfTruth == externalSystem.ExternalSystemId {
		if e.Channel == fields.Channel &&
			e.ChannelData == fields.ChannelData &&
			e.Content == fields.Content &&
			e.ContentType == fields.ContentType &&
			e.Identifier == fields.Identifier &&
			e.EventType == fields.EventType {
			return true
		}
	}

	return false
}

type Sender struct {
	Participant  common.Participant `json:"participant"`
	RelationType string             `json:"relationType"`
}

type Receiver struct {
	Participant  common.Participant `json:"participant"`
	RelationType string             `json:"relationType"`
}

func (s Sender) Available() bool {
	return s.Participant.ID != "" && s.Participant.ParticipantType != ""
}

func (r Receiver) Available() bool {
	return r.Participant.ID != "" && r.Participant.ParticipantType != ""
}

func (e *InteractionEvent) HasExternalSystem(externalSystem common.ExternalSystem) bool {
	for _, es := range e.ExternalSystems {
		if es.ExternalSystemId == externalSystem.ExternalSystemId &&
			es.ExternalId == externalSystem.ExternalId &&
			es.ExternalSource == externalSystem.ExternalSource &&
			es.ExternalUrl == externalSystem.ExternalUrl &&
			es.ExternalIdSecond == externalSystem.ExternalIdSecond {
			return true
		}
	}
	return false
}
