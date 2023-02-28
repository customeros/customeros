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
}

func (t TicketData) HasCollaborators() bool {
	return len(t.CollaboratorUserExternalIds) > 0
}
