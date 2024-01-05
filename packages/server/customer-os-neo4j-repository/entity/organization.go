package entity

import (
	"time"
)

type OrganizationEntity struct {
	ID                 string
	CustomerOsId       string
	Name               string
	Description        string
	Website            string
	Industry           string
	SubIndustry        string
	IndustryGroup      string
	TargetAudience     string
	ValueProposition   string
	IsPublic           bool
	IsCustomer         bool
	Hide               bool
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	ReferenceId        string
	Note               string
	Employees          int64
	CreatedAt          time.Time
	LastTouchpointAt   *time.Time
	UpdatedAt          time.Time
	LastTouchpointId   *string
	Source             DataSource
	SourceOfTruth      DataSource
	AppSource          string
	YearFounded        *int64
	Headquarters       string
	EmployeeGrowthRate string
	LogoUrl            string

	LinkedOrganizationType *string

	SuggestedMerge struct {
		SuggestedAt *time.Time
		SuggestedBy *string
		Confidence  *float64
	}
	RenewalSummary    RenewalSummary
	OnboardingDetails OnboardingDetails
	WebScrapeDetails  WebScrapeDetails
}

type RenewalSummary struct {
	ArrForecast            *float64
	MaxArrForecast         *float64
	NextRenewalAt          *time.Time
	RenewalLikelihood      string
	RenewalLikelihoodOrder *int64
}

type OnboardingDetails struct {
	Status       string
	SortingOrder *int64
	UpdatedAt    *time.Time
	Comments     string
}

type WebScrapeDetails struct {
	WebScrapedUrl             string
	WebScrapedAt              *time.Time
	WebScrapeLastRequestedAt  *time.Time
	WebScrapeLastRequestedUrl string
	WebScrapeAttempts         int64
}

type OrganizationEntities []OrganizationEntity
