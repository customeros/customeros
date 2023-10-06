package entity

import (
	"time"
)

type OrganizationEntity struct {
	ID                string
	CustomerOsId      string
	ReferenceId       string
	Name              string
	Description       string
	Website           string
	Industry          string
	SubIndustry       string
	IndustryGroup     string
	TargetAudience    string
	ValueProposition  string
	IsPublic          bool
	Hide              bool
	Market            string
	LastFundingRound  string
	LastFundingAmount string
	Note              string
	Employees         int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	LastTouchpointAt  *time.Time
	LastTouchpointId  *string
	Source            DataSource
	SourceOfTruth     DataSource
	AppSource         string

	LinkedOrganizationType *string

	SuggestedMerge struct {
		SuggestedAt *time.Time
		SuggestedBy *string
		Confidence  *float64
	}
	RenewalLikelihood RenewalLikelihood
	RenewalForecast   RenewalForecast
	BillingDetails    BillingDetails
}

type RenewalLikelihood struct {
	RenewalLikelihood         string `neo4jDb:"property:renewalLikelihood;lookupName:RENEWAL_LIKELIHOOD;supportCaseSensitive:false"`
	PreviousRenewalLikelihood string
	Comment                   *string
	UpdatedAt                 *time.Time
	UpdatedBy                 *string
}

type RenewalForecast struct {
	Amount          *float64 `neo4jDb:"property:renewalForecastAmount;lookupName:FORECAST_AMOUNT;supportCaseSensitive:false"`
	PotentialAmount *float64
	Comment         *string
	UpdatedAt       *time.Time
	UpdatedById     *string
}

type BillingDetails struct {
	Amount            *float64
	Frequency         string
	RenewalCycle      string
	RenewalCycleStart *time.Time
	RenewalCycleNext  *time.Time `neo4jDb:"property:billingDetailsRenewalCycleNext;lookupName:RENEWAL_CYCLE_NEXT;supportCaseSensitive:false"`
}

type OrganizationEntities []OrganizationEntity
