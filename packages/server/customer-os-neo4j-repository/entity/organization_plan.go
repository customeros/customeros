package entity

import (
	"time"
)

type OrganizationPlanStatusDetails struct {
	Status    string
	UpdatedAt time.Time
	Comments  string
}

type OrganizationPlanEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Retired       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	StatusDetails OrganizationPlanStatusDetails
	MasterPlanId  string
}

type OrganizationPlanEntities []OrganizationPlanEntity
