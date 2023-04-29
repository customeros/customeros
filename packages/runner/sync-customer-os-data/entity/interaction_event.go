package entity

import "time"

type ParticipantType int

const (
	ORGANIZATION ParticipantType = 1
	USER         ParticipantType = 2
	CONTACT      ParticipantType = 3
	EMAIL        ParticipantType = 4
	PHONE        ParticipantType = 5
)

type InteractionEventParticipant struct {
	ExternalId      string
	ParticipantType ParticipantType
	RelationType    string
}

func (participant InteractionEventParticipant) GetNodeLabel() string {
	switch participant.ParticipantType {
	case ORGANIZATION:
		return "Organization"
	case USER:
		return "User"
	case CONTACT:
		return "Contact"
	case EMAIL:
		return "Email"
	case PHONE:
		return "PhoneNumber"
	default:
		return "Unknown"
	}
}

type InteractionEventData struct {
	Id               string
	Content          string
	ContentType      string
	Type             string
	Channel          string
	CreatedAt        time.Time
	PartOfExternalId string
	SentBy           InteractionEventParticipant
	SentTo           map[string]InteractionEventParticipant

	ExternalId     string
	ExternalSyncId string
	ExternalSystem string
}

func (data InteractionEventData) IsPartOf() bool {
	return len(data.PartOfExternalId) > 0
}

func (data InteractionEventData) HasSender() bool {
	return len(data.SentBy.ExternalId) > 0
}

func (data InteractionEventData) HasRecipients() bool {
	return len(data.SentTo) > 0
}
