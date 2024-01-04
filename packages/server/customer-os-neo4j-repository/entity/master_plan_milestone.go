package entity

import (
	"time"
)

type MasterPlanMilestoneEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
	Retired       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type MasterPlanMilestoneEntities []MasterPlanMilestoneEntity
