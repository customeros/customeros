package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type OrganizationDataFields struct {
	Name              string
	Hide              bool
	Description       string
	Website           string
	Industry          string
	SubIndustry       string
	IndustryGroup     string
	TargetAudience    string
	ValueProposition  string
	IsPublic          bool
	Employees         int64
	Market            string
	LastFundingRound  string
	LastFundingAmount string
}

type OrganizationFields struct {
	ID                     string
	Tenant                 string
	IgnoreEmptyFields      bool
	OrganizationDataFields OrganizationDataFields
	Source                 commonModels.Source
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
	RenewalLikelihood      *RenewalLikelihoodFields
}

type RenewalLikelihoodFields struct {
	RenewalLikelihood RenewalLikelihoodProbability
	Comment           *string
	UpdatedAt         time.Time
	UpdatedBy         string `validate:"required"`
}

type RenewalForecastFields struct {
	Amount          *float64
	PotentialAmount *float64
	Comment         *string
	UpdatedAt       time.Time
	UpdatedBy       string
}

type BillingDetailsFields struct {
	Amount            *float64
	Frequency         string
	RenewalCycle      string
	RenewalCycleStart *time.Time
	RenewalCycleNext  *time.Time
	UpdatedBy         string
}
