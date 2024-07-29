package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"
)

type Issue struct {
	ID                        string                  `json:"id"`
	Tenant                    string                  `json:"tenant"`
	Subject                   string                  `json:"subject"`
	Description               string                  `json:"description"`
	Status                    string                  `json:"status"`
	Priority                  string                  `json:"priority"`
	ReportedByOrganizationId  string                  `json:"reportedByOrganizationId,omitempty"`
	SubmittedByOrganizationId string                  `json:"submittedByOrganizationId,omitempty"`
	SubmittedByUserId         string                  `json:"submittedByUserId,omitempty"`
	Source                    cmnmod.Source           `json:"source"`
	ExternalSystems           []cmnmod.ExternalSystem `json:"externalSystem"`
	CreatedAt                 time.Time               `json:"createdAt,omitempty"`
	UpdatedAt                 time.Time               `json:"updatedAt,omitempty"`
	AssignedToUserIds         []string                `json:"assignedToUserIds,omitempty"`
	FollowedByUserIds         []string                `json:"followedByUserIds,omitempty"`
}

func (i *Issue) AddAssignedToUserId(userId string) {
	i.AssignedToUserIds = utils.AddToListIfNotExists(i.AssignedToUserIds, userId)
}

func (i *Issue) RemoveAssignedToUserId(userId string) {
	i.AssignedToUserIds = utils.RemoveFromList(i.AssignedToUserIds, userId)
}

func (i *Issue) AddFollowedByUserId(userId string) {
	i.FollowedByUserIds = utils.AddToListIfNotExists(i.FollowedByUserIds, userId)
}

func (i *Issue) RemoveFollowedByUserId(userId string) {
	i.FollowedByUserIds = utils.RemoveFromList(i.FollowedByUserIds, userId)
}

func (i *Issue) SameData(fields IssueDataFields, externalSystem cmnmod.ExternalSystem) bool {
	if !externalSystem.Available() {
		return false
	}

	if externalSystem.Available() && !i.HasExternalSystem(externalSystem) {
		return false
	}

	if i.Subject == fields.Subject &&
		i.Description == fields.Description &&
		i.Status == fields.Status &&
		i.Priority == fields.Priority &&
		i.ReportedByOrganizationId == utils.IfNotNilString(fields.ReportedByOrganizationId) &&
		i.SubmittedByOrganizationId == utils.IfNotNilString(fields.SubmittedByOrganizationId) &&
		i.SubmittedByUserId == utils.IfNotNilString(fields.SubmittedByUserId) {
		return true
	}

	return false
}

func (i *Issue) HasExternalSystem(externalSystem cmnmod.ExternalSystem) bool {
	for _, es := range i.ExternalSystems {
		if es.ExternalSystemId == externalSystem.ExternalSystemId &&
			es.ExternalId == externalSystem.ExternalId &&
			es.ExternalSource == externalSystem.ExternalSource &&
			es.ExternalUrl == externalSystem.ExternalUrl &&
			es.ExternalIdSecond == externalSystem.ExternalIdSecond {
			return true
		}
	}
	return false
}
