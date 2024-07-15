package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"time"
)

type MasterPlan struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	Retired      bool                           `json:"retired"`
	CreatedAt    time.Time                      `json:"createdAt"`
	UpdatedAt    time.Time                      `json:"updatedAt"`
	SourceFields common.Source                  `json:"source"`
	Milestones   map[string]MasterPlanMilestone `json:"milestones"`
}

type MasterPlanMilestone struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Retired       bool          `json:"retired"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	SourceFields  common.Source `json:"source"`
	Optional      bool          `json:"optional"`
	Order         int64         `json:"order"`
	DurationHours int64         `json:"durationHours"`
	Items         []string      `json:"items"`
}
