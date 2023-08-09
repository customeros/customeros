package entity

import (
	"strings"
)

func (participant InteractionEventParticipant) GetNodeLabel() string {
	switch strings.ToUpper(participant.ParticipantType) {
	case "ORGANIZATION":
		return "Organization"
	case "USER":
		return "User"
	case "CONTACT":
		return "Contact"
	case "EMAIL":
		return "Email"
	case "PHONE":
		return "PhoneNumber"
	default:
		return "Unknown"
	}
}

type InteractionEventParticipant struct {
	ExternalId      string `json:"externalId,omitempty"`
	ParticipantType string `json:"participantType,omitempty"`
	RelationType    string `json:"relationType,omitempty"`
}

type InteractionEventData struct {
	BaseData
	Content          string                        `json:"content,omitempty"`
	ContentType      string                        `json:"contentType,omitempty"`
	Type             string                        `json:"type,omitempty"`
	Channel          string                        `json:"channel,omitempty"`
	PartOfExternalId string                        `json:"partOfExternalId,omitempty"`
	SentBy           InteractionEventParticipant   `json:"sentBy,omitempty"`
	SentTo           []InteractionEventParticipant `json:"sentTo,omitempty"`
}

func (i *InteractionEventData) IsPartOf() bool {
	return len(i.PartOfExternalId) > 0
}

func (i *InteractionEventData) HasSender() bool {
	return len(i.SentBy.ExternalId) > 0
}

func (i *InteractionEventData) HasRecipients() bool {
	return len(i.SentTo) > 0
}

func (i *InteractionEventData) Normalize() {
	i.SetTimes()
}
