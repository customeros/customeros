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
	FollowerUserExternalIds     []string
	SubmitterExternalId         string
	RequesterExternalId         string
	AssigneeUserExternalId      string
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

func (t TicketData) HasFollowers() bool {
	return len(t.FollowerUserExternalIds) > 0
}

func (t TicketData) HasAssignee() bool {
	return len(t.AssigneeUserExternalId) > 0
}
