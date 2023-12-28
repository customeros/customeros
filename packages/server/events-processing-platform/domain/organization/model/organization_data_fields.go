package model

import (
	"time"

	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

type OrganizationDataFields struct {
	Name               string
	Hide               bool
	Description        string
	Website            string
	Industry           string
	SubIndustry        string
	IndustryGroup      string
	TargetAudience     string
	ValueProposition   string
	IsPublic           bool
	IsCustomer         bool
	Employees          int64
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	ReferenceId        string
	Note               string
	YearFounded        *int64
	Headquarters       string
	EmployeeGrowthRate string
	LogoUrl            string
}

type OrganizationFields struct {
	ID                     string
	Tenant                 string
	OrganizationDataFields OrganizationDataFields
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
	Source                 cmnmod.Source
	ExternalSystem         cmnmod.ExternalSystem
}
