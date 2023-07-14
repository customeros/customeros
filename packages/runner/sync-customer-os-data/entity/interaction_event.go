package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

/*
{
  "content": "Hello World!",
  "contentType": "text/plain",
  "type": "message",
  "channel": "email",
  "partOfExternalId": "conversation-123",

  "sentBy": {
    "externalId": "user-123",
    "participantType": "person",
    "relationType": "from"
  },

  "sentTo": {
    "contact-456": {
      "externalId": "contact-456",
      "participantType": "person",
      "relationType": "to"
    },
    "group-789": {
      "externalId": "group-789",
      "participantType": "group",
      "relationType": "to"
    }
  },

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "event-123",
  "externalSystem": "HubSpot",
  "createdAt": "2023-03-05T17:11:22Z",
  "updatedAt": "2023-03-05T17:11:22Z",
  "syncId": "sync-1234"
}
*/

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
	Content          string                                 `json:"content,omitempty"`
	ContentType      string                                 `json:"contentType,omitempty"`
	Type             string                                 `json:"type,omitempty"`
	Channel          string                                 `json:"channel,omitempty"`
	PartOfExternalId string                                 `json:"partOfExternalId,omitempty"`
	SentBy           InteractionEventParticipant            `json:"sentBy,omitempty"`
	SentTo           map[string]InteractionEventParticipant `json:"sentTo,omitempty"`
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

func (i *InteractionEventData) FormatTimes() {
	if i.CreatedAt != nil {
		i.CreatedAt = common_utils.TimePtr((*i.CreatedAt).UTC())
	} else {
		i.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if i.UpdatedAt != nil {
		i.UpdatedAt = common_utils.TimePtr((*i.UpdatedAt).UTC())
	} else {
		i.UpdatedAt = common_utils.TimePtr(common_utils.Now())
	}
}

func (i *InteractionEventData) Normalize() {
	i.FormatTimes()
}
