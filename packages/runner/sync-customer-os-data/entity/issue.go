package entity

import (
	utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type IssueData struct {
	BaseData
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

func (i *IssueData) Normalize() {
	i.SetTimes()
	i.Tags = utils.FilterOutEmpty(i.Tags)
	i.Tags = utils.LowercaseSliceOfStrings(i.Tags)
	i.Tags = utils.RemoveDuplicates(i.Tags)
}
