package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type InteractionEventParticipant struct {
	ReferencedUser         ReferencedUser         `json:"referencedUser,omitempty"`
	ReferencedContact      ReferencedContact      `json:"referencedContact,omitempty"`
	ReferencedOrganization ReferencedOrganization `json:"referencedOrganization,omitempty"`
	ReferencedParticipant  ReferencedParticipant  `json:"referencedParticipant,omitempty"`
	ReferencedJobRole      ReferencedJobRole      `json:"referencedJobRole,omitempty"`
	RelationType           string                 `json:"relationType,omitempty"`
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
	// in sent to or sent by at least 1 contact should be available in the system
	ContactRequired bool `json:"contactRequired,omitempty"`
}

func (i *InteractionEventData) IsPartOfByExternalId() bool {
	return len(i.PartOfExternalId) > 0
}

func (i *InteractionEventData) HasSender() bool {
	return i.SentBy.ReferencedUser.Available() ||
		i.SentBy.ReferencedContact.Available() ||
		i.SentBy.ReferencedOrganization.Available() ||
		i.SentBy.ReferencedParticipant.Available() ||
		i.SentBy.ReferencedJobRole.Available()
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
