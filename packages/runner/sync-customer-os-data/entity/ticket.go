package entity

import "time"

type IssueData struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ExternalId     string
	ExternalSystem string
	ExternalSyncId string
	ExternalUrl    string

	Subject                        string
	Status                         string
	Priority                       string
	Description                    string
	Tags                           []string
	CollaboratorUserExternalIds    []string
	FollowerUserExternalIds        []string
	ReporterOrganizationExternalId string
	AssigneeUserExternalId         string
}

func (t IssueData) HasCollaboratorUsers() bool {
	return len(t.CollaboratorUserExternalIds) > 0
}

func (t IssueData) HasReporterOrganization() bool {
	return len(t.ReporterOrganizationExternalId) > 0
}

func (t IssueData) HasFollowerUsers() bool {
	return len(t.FollowerUserExternalIds) > 0
}

func (t IssueData) HasAssignee() bool {
	return len(t.AssigneeUserExternalId) > 0
}

func (t IssueData) HasTags() bool {
	return len(t.Tags) > 0
}
