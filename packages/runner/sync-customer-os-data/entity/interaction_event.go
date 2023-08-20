package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
		return ""
	}
}

type InteractionEventParticipant struct {
	OpenlineId                string `json:"openlineId,omitempty"`
	ExternalId                string `json:"externalId,omitempty"`
	ParticipantType           string `json:"participantType,omitempty"`
	RelationType              string `json:"relationType,omitempty"`
	ReplaceContactWithJobRole bool   `json:"replaceContactWithJobRole,omitempty"`
	OrganizationId            string `json:"organizationId,omitempty"`
}

type InteractionEventData struct {
	BaseData
	Content          string                        `json:"content,omitempty"`
	ContentType      string                        `json:"contentType,omitempty"`
	Type             string                        `json:"type,omitempty"`
	Channel          string                        `json:"channel,omitempty"`
	Identifier       string                        `json:"identifier,omitempty"`
	Hide             bool                          `json:"hide,omitempty"`
	PartOfExternalId string                        `json:"partOfExternalId,omitempty"`
	PartOfSession    InteractionSession            `json:"partOfSession,omitempty"`
	SentBy           InteractionEventParticipant   `json:"sentBy,omitempty"`
	SentTo           []InteractionEventParticipant `json:"sentTo,omitempty"`
}

func (i *InteractionEventData) IsPartOfByExternalId() bool {
	return len(i.PartOfExternalId) > 0
}

func (i *InteractionEventData) HasSender() bool {
	return len(i.SentBy.ExternalId) > 0
}

func (i *InteractionEventData) HasSession() bool {
	return i.PartOfSession.ExternalId != ""
}

func (i *InteractionEventData) Normalize() {
	i.SetTimes()
	if i.HasSession() {
		if i.PartOfSession.CreatedAtStr != "" && i.PartOfSession.CreatedAt == nil {
			i.PartOfSession.CreatedAt, _ = utils.UnmarshalDateTime(i.PartOfSession.CreatedAtStr)
		}
		if i.PartOfSession.CreatedAt != nil {
			i.PartOfSession.CreatedAt = common_utils.TimePtr((*i.PartOfSession.CreatedAt).UTC())
		} else {
			i.PartOfSession.CreatedAt = common_utils.TimePtr(common_utils.Now())
		}
	}

}
