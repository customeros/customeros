package entity

import "time"

type TicketData struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ExternalId     string
	ExternalSystem string
	ExternalSyncId string
	ExternalUrl    string

	Subject                     string
	CollaboratorUserExternalIds []string
	SubmitterExternalId         string
	RequesterExternalId         string
}

func (t TicketData) HasCollaborators() bool {
	return len(t.CollaboratorUserExternalIds) > 0
}

func (t TicketData) HasSubmitter() bool {
	return len(t.SubmitterExternalId) > 0
}

func (t TicketData) HasRequester() bool {
	return len(t.RequesterExternalId) > 0
}
