package entity

import (
	"time"
)

type MasterPlanEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Name          string
}
