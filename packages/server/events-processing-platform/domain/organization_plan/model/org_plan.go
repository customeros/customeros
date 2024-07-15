package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"time"
)

type OrganizationPlanDetails struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Comments  string    `json:"comments"`
}

type OrganizationPlanMilestoneItem struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      string    `json:"text"`
	Uuid      string    `json:"uuid"`
}

type OrganizationPlan struct {
	ID            string                               `json:"id"`
	Name          string                               `json:"name"`
	Retired       bool                                 `json:"retired"`
	CreatedAt     time.Time                            `json:"createdAt"`
	UpdatedAt     time.Time                            `json:"updatedAt"`
	SourceFields  common.Source                        `json:"source"`
	Milestones    map[string]OrganizationPlanMilestone `json:"milestones"`
	StatusDetails OrganizationPlanDetails              `json:"statusDetails"`
	MasterPlanId  string                               `json:"masterPlanId"`
}

type OrganizationPlanMilestone struct {
	ID            string                          `json:"id"`
	Name          string                          `json:"name"`
	Retired       bool                            `json:"retired"`
	CreatedAt     time.Time                       `json:"createdAt"`
	UpdatedAt     time.Time                       `json:"updatedAt"`
	SourceFields  common.Source                   `json:"source"`
	Optional      bool                            `json:"optional"`
	Order         int64                           `json:"order"`
	DueDate       time.Time                       `json:"dueDate"`
	Items         []OrganizationPlanMilestoneItem `json:"items"`
	StatusDetails OrganizationPlanDetails         `json:"statusDetails"`
	Adhoc         bool                            `json:"adhoc"`
}
