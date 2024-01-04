package entity

import (
	"time"
)

type MasterPlanEntity struct {
	SourceFields
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Retired   bool
}

type MasterPlanEntities []MasterPlanEntity
