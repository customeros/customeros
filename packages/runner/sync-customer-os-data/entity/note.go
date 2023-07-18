package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type NoteData struct {
	BaseData
	Html                          string   `json:"html,omitempty"`
	Text                          string   `json:"text,omitempty"`
	CreatorUserExternalId         string   `json:"externalUserId,omitempty"`
	CreatorUserExternalOwnerId    string   `json:"externalOwnerId,omitempty"`
	CreatorExternalId             string   `json:"externalCreatorId,omitempty"`
	NotedContactsExternalIds      []string `json:"contactsExternalIds,omitempty"`
	NotedOrganizationsExternalIds []string `json:"organizationsExternalIds,omitempty"`
	MentionedTags                 []string `json:"mentionedTags,omitempty"`
	MentionedIssueExternalId      string   `json:"mentionedIssueExternalId,omitempty"`
}

func (n *NoteData) HasNotedContacts() bool {
	return len(n.NotedContactsExternalIds) > 0
}

func (n *NoteData) HasNotedOrganizations() bool {
	return len(n.NotedOrganizationsExternalIds) > 0
}

func (n *NoteData) HasMentionedTags() bool {
	return len(n.MentionedTags) > 0
}

func (n *NoteData) HasCreatorUser() bool {
	return len(n.CreatorUserExternalId) > 0
}

func (n *NoteData) HasCreatorUserOwner() bool {
	return len(n.CreatorUserExternalOwnerId) > 0
}

func (n *NoteData) HasCreator() bool {
	return len(n.CreatorExternalId) > 0
}

func (n *NoteData) HasMentionedIssue() bool {
	return n.MentionedIssueExternalId != ""
}

func (n *NoteData) FormatTimes() {
	if n.CreatedAt != nil {
		n.CreatedAt = common_utils.TimePtr((*n.CreatedAt).UTC())
	} else {
		n.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
}

func (n *NoteData) Normalize() {
	n.FormatTimes()
}
