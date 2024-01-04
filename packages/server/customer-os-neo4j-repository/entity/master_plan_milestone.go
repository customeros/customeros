package entity

import (
	"time"
)

type MasterPlanMilestoneEntity struct {
	DataLoaderKey
	SourceFields
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
	Retired       bool
}

type MasterPlanMilestoneEntities []MasterPlanMilestoneEntity
