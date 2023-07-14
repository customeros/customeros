package entity

import (
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

/*
{
  "html": "<b>Hello</b> world!",
  "text": "Hello world!",
  "externalUserId": "user-123",
  "externalOwnerId": "owner-456",
  "externalCreatorId": "system-xyz",
  "contactsExternalIds": [
    "contact-123",
    "contact-456"
  ],
  "organizationsExternalIds": [
    "org-123"
  ],
  "mentionedTags": [
    "important",
    "update"
  ],
  "mentionedIssueExternalId": "issue-123",

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "abcd1234",
  "externalSystem": "HubSpot",
  "createdAt": "2022-02-28T19:52:05Z",
  "updatedAt": "2022-03-01T11:23:45Z",
  "syncId": "sync_1234"
}
*/

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
