package model

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

type OrgPlanDetails struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Comments  string    `json:"comments"`
}

type OrgPlanMilestoneTask struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      string    `json:"text"`
}

type OrgPlan struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Retired      bool                        `json:"retired"`
	CreatedAt    time.Time                   `json:"createdAt"`
	UpdatedAt    time.Time                   `json:"updatedAt"`
	SourceFields commonmodel.Source          `json:"source"`
	Milestones   map[string]OrgPlanMilestone `json:"milestones"`
	Details      OrgPlanDetails              `json:"details"`
	MasterPlanId string                      `json:"masterPlanId"`
}

type OrgPlanMilestone struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Retired       bool                   `json:"retired"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	SourceFields  commonmodel.Source     `json:"source"`
	Optional      bool                   `json:"optional"`
	Order         int64                  `json:"order"`
	DurationHours int64                  `json:"durationHours"`
	Items         []OrgPlanMilestoneTask `json:"items"` // maybe make this a map of `item text -> item done` to keep this very simple
	Details       OrgPlanDetails         `json:"details"`
}
