package entity

import (
	utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

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
