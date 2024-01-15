package entity

import (
	"time"
)

type OrganizationPlanMilestoneItem struct {
	Text      string
	Status    string
	UpdatedAt time.Time
}

type OrganizationPlanMilestoneStatusDetails struct {
	Status    string
	UpdatedAt time.Time
	Comments  string
}

type OrganizationPlanMilestoneEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Order         int64
	DurationHours int64
	Items         []OrganizationPlanMilestoneItem
	Optional      bool
	Retired       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	StatusDetails OrganizationPlanMilestoneStatusDetails
}

type OrganizationPlanMilestoneEntities []OrganizationPlanMilestoneEntity
