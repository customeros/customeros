package entity

import "time"

type NoteData struct {
	Id                            string
	Html                          string
	Text                          string
	CreatedAt                     time.Time
	CreatorUserExternalId         string
	CreatorUserExternalOwnerId    string
	CreatorExternalId             string
	NotedContactsExternalIds      []string
	NotedOrganizationsExternalIds []string
	MentionedTags                 []string
	ExternalId                    string
	ExternalSyncId                string
	ExternalSystem                string
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
