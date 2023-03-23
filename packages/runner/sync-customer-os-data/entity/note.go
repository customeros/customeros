package entity

import "time"

type NoteData struct {
	Id                             string
	Html                           string
	Text                           string
	CreatedAt                      time.Time
	CreatorUserExternalId          string
	CreatorUserExternalOwnerId     string
	CreatorUserOrContactExternalId string
	NotedContactsExternalIds       []string
	NotedOrganizationsExternalIds  []string
	MentionedIssuesExternalIds     []string
	ExternalId                     string
	ExternalSyncId                 string
	ExternalSystem                 string
}

func (n NoteData) HasNotedContacts() bool {
	return len(n.NotedContactsExternalIds) > 0
}

func (n NoteData) HasNotedOrganizations() bool {
	return len(n.NotedOrganizationsExternalIds) > 0
}

func (n NoteData) HasMentionedIssues() bool {
	return len(n.MentionedIssuesExternalIds) > 0
}

func (n NoteData) HasCreatorUser() bool {
	return len(n.CreatorUserExternalId) > 0
}

func (n NoteData) HasCreatorUserOwner() bool {
	return len(n.CreatorUserExternalOwnerId) > 0
}

func (n NoteData) HasCreatorUserOrContact() bool {
	return len(n.CreatorUserOrContactExternalId) > 0
}
