package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Issue struct {
	ID                      string                  `json:"id"`
	Tenant                  string                  `json:"tenant"`
	Subject                 string                  `json:"subject"`
	Description             string                  `json:"description"`
	Status                  string                  `json:"status"`
	Priority                string                  `json:"priority"`
	ReportedByOrganization  string                  `json:"reportedByOrganization,omitempty"`
	SubmittedByOrganization string                  `json:"submittedByOrganization,omitempty"`
	SubmittedByUser         string                  `json:"submittedByUser,omitempty"`
	Source                  cmnmod.Source           `json:"source"`
	ExternalSystems         []cmnmod.ExternalSystem `json:"externalSystem"`
	CreatedAt               time.Time               `json:"createdAt,omitempty"`
	UpdatedAt               time.Time               `json:"updatedAt,omitempty"`
	AssignedToUserIds       []string                `json:"assignedToUserIds,omitempty"`
	FollowedByUserIds       []string                `json:"followedByUserIds,omitempty"`
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
