package entity

import (
	"time"
)

type OrganizationPlanMilestoneItem struct {
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Uuid      string    `json:"uuid"`
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
	DueDate       time.Time
	Items         []OrganizationPlanMilestoneItem
	Optional      bool
	Retired       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	StatusDetails OrganizationPlanMilestoneStatusDetails
	Adhoc         bool
}

type OrganizationPlanMilestoneEntities []OrganizationPlanMilestoneEntity
