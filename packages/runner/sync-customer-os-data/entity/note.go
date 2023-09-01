package entity

type NoteData struct {
	BaseData
	Content                       string   `json:"content,omitempty"`
	ContentType                   string   `json:"contentType,omitempty"`
	Text                          string   `json:"text,omitempty"`
	CreatorUserExternalId         string   `json:"externalUserId,omitempty"`
	CreatorUserExternalOwnerId    string   `json:"externalOwnerId,omitempty"`
	CreatorExternalId             string   `json:"externalCreatorId,omitempty"`
	NotedContactsExternalIds      []string `json:"externalContactsIds,omitempty"`
	NotedOrganizationsExternalIds []string `json:"externalOrganizationsIds,omitempty"`
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

func (n *NoteData) Normalize() {
	n.SetTimes()
}
