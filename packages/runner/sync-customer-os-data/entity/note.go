package entity

import "time"

type NoteData struct {
	Id                            string    `json:"id,omitempty"`
	Html                          string    `json:"html,omitempty"`
	Text                          string    `json:"text,omitempty"`
	CreatedAt                     time.Time `json:"createdAt,omitempty"`
	CreatorUserExternalId         string    `json:"externalUserId,omitempty"`
	CreatorUserExternalOwnerId    string    `json:"externalOwnerId,omitempty"`
	CreatorExternalId             string    `json:"externalCreatorId,omitempty"`
	NotedContactsExternalIds      []string  `json:"contactsExternalIds,omitempty"`
	NotedOrganizationsExternalIds []string  `json:"organizationsExternalIds,omitempty"`
	MentionedTags                 []string  `json:"mentionedTags,omitempty"`
	ExternalId                    string    `json:"externalId,omitempty"`
	ExternalSyncId                string    `json:"externalSyncId,omitempty"`
	ExternalSystem                string    `json:"externalSystem,omitempty"`
}

func (n NoteData) HasNotedContacts() bool {
	return len(n.NotedContactsExternalIds) > 0
}

func (n NoteData) HasNotedOrganizations() bool {
	return len(n.NotedOrganizationsExternalIds) > 0
}

func (n NoteData) HasMentionedTags() bool {
	return len(n.MentionedTags) > 0
}

func (n NoteData) HasCreatorUser() bool {
	return len(n.CreatorUserExternalId) > 0
}

func (n NoteData) HasCreatorUserOwner() bool {
	return len(n.CreatorUserExternalOwnerId) > 0
}

func (n NoteData) HasCreator() bool {
	return len(n.CreatorExternalId) > 0
}
