package entity

import (
	utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

/*
{
  "externalUrl": "https://issue-tracker.com/issue/123",
  "subject": "Fix login bug",
  "status": "Open",
  "priority": "High",
  "description": "The login page is throwing an error when submitting credentials",

  "tags": [
    "bug",
    "critical",
    "login"
  ],

  "collaboratorUserExternalIds": [
    "user-1",
    "user-2"
  ],

  "followerUserExternalIds": [
    "user-3",
    "user-4"
  ],

  "reporterOrganizationExternalId": "org-abc",
  "assigneeUserExternalId": "user-5",

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "issue-123",
  "externalSystem": "Jira",
  "createdAt": "2023-03-01T12:34:56Z",
  "updatedAt": "2023-03-02T15:19:00Z",
  "syncId": "sync-1234"
}
*/

type IssueData struct {
	BaseData

	ExternalUrl string   `json:"externalUrl,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	Status      string   `json:"status,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	CollaboratorUserExternalIds []string `json:"collaboratorUserExternalIds,omitempty"`
	FollowerUserExternalIds     []string `json:"followerUserExternalIds,omitempty"`

	ReporterOrganizationExternalId string `json:"reporterOrganizationExternalId,omitempty"`
	AssigneeUserExternalId         string `json:"assigneeUserExternalId,omitempty"`
}

func (i *IssueData) HasCollaboratorUsers() bool {
	return len(i.CollaboratorUserExternalIds) > 0
}

func (i *IssueData) HasReporterOrganization() bool {
	return len(i.ReporterOrganizationExternalId) > 0
}

func (i *IssueData) HasFollowerUsers() bool {
	return len(i.FollowerUserExternalIds) > 0
}

func (i *IssueData) HasAssignee() bool {
	return len(i.AssigneeUserExternalId) > 0
}

func (i *IssueData) HasTags() bool {
	return len(i.Tags) > 0
}

func (i *IssueData) FormatTimes() {
	if i.CreatedAt != nil {
		i.CreatedAt = utils.TimePtr((*i.CreatedAt).UTC())
	} else {
		i.CreatedAt = utils.TimePtr(utils.Now())
	}
	if i.UpdatedAt != nil {
		i.UpdatedAt = utils.TimePtr((*i.UpdatedAt).UTC())
	} else {
		i.UpdatedAt = utils.TimePtr(utils.Now())
	}
}

func (i *IssueData) Normalize() {
	i.FormatTimes()
	utils.FilterEmpty(i.Tags)
	utils.LowercaseStrings(i.Tags)
	i.Tags = utils.RemoveDuplicates(i.Tags)
}
