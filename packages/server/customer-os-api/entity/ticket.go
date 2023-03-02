package entity

import "time"

type Ticket struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ExternalId     string
	ExternalSystem string
	ExternalSyncId string
	ExternalUrl    string

	Subject                     string
	Status                      string
	Priority                    string
	Description                 string
	Tags                        []string
	CollaboratorUserExternalIds []string
	FollowerUserExternalIds     []string
	SubmitterExternalId         string
	RequesterExternalId         string
	AssigneeUserExternalId      string
}
