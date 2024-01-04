package entity

import (
	"time"
)

type MasterPlanMilestoneEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
}

type MasterPlanMilestoneEntities []MasterPlanMilestoneEntity
