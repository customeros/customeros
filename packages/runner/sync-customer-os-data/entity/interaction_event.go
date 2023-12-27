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

type BelongsTo struct {
	Issue   ReferencedIssue              `json:"issue,omitempty"`
	Session ReferencedInteractionSession `json:"session,omitempty"`
}

func (b BelongsTo) Available() bool {
	return b.Issue.Available() || b.Session.Available()
}

type InteractionEventData struct {
	BaseData
	Content        string                        `json:"content,omitempty"`
	ContentType    string                        `json:"contentType,omitempty"`
	EventType      string                        `json:"eventType,omitempty"`
	Channel        string                        `json:"channel,omitempty"`
	Identifier     string                        `json:"identifier,omitempty"`
	Hide           bool                          `json:"hide,omitempty"`
	BelongsTo      BelongsTo                     `json:"belongsTo,omitempty"`
	SessionDetails InteractionSession            `json:"sessionDetails,omitempty"`
	SentBy         InteractionEventParticipant   `json:"sentBy,omitempty"`
	SentTo         []InteractionEventParticipant `json:"sentTo,omitempty"`
	// in sent to or sent by at least 1 contact should be available in the system
	ContactRequired bool `json:"contactRequired,omitempty"`
	// interaction session should already exist in the system
	ParentRequired bool `json:"parentRequired,omitempty"`
}

func (i *InteractionEventData) IsPartOf() bool {
	return i.BelongsTo.Available()
}

func (i *InteractionEventData) HasSender() bool {
	return i.SentBy.ReferencedUser.Available() ||
		i.SentBy.ReferencedContact.Available() ||
		i.SentBy.ReferencedOrganization.Available() ||
		i.SentBy.ReferencedParticipant.Available() ||
		i.SentBy.ReferencedJobRole.Available()
}

func (i *InteractionEventData) HasSession() bool {
	return i.SessionDetails.ExternalId != ""
}

func (i *InteractionEventData) Normalize() {
	i.SetTimes()
	if i.HasSession() {
		if i.SessionDetails.CreatedAtStr != "" && i.SessionDetails.CreatedAt == nil {
			i.SessionDetails.CreatedAt, _ = utils.UnmarshalDateTime(i.SessionDetails.CreatedAtStr)
		}
		if i.SessionDetails.CreatedAt != nil {
			i.SessionDetails.CreatedAt = common_utils.TimePtr((*i.SessionDetails.CreatedAt).UTC())
		} else {
			i.SessionDetails.CreatedAt = common_utils.TimePtr(common_utils.Now())
		}
	}

}
