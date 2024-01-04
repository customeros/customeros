package entity

import (
	"time"
)

type MasterPlanEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Retired       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type MasterPlanEntities []MasterPlanEntity
