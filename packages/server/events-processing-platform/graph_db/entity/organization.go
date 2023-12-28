package entity

import (
	"fmt"
	"time"
)

type OnboardingStatus string

const (
	OnboardingStatusNotApplicable OnboardingStatus = "NOT_APPLICABLE"
	OnboardingStatusNotStarted    OnboardingStatus = "NOT_STARTED"
	OnboardingStatusOnTrack       OnboardingStatus = "ON_TRACK"
	OnboardingStatusLate          OnboardingStatus = "LATE"
	OnboardingStatusStuck         OnboardingStatus = "STUCK"
	OnboardingStatusDone          OnboardingStatus = "DONE"
	OnboardingStatusSuccessful    OnboardingStatus = "SUCCESSFUL"
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

	InteractionEventParticipantDetails InteractionEventParticipantDetails
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

func (organization OrganizationEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organization.ID, organization.Name)
}

func (OrganizationEntity) IsNotedEntity() {}

func (OrganizationEntity) NotedEntityLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsInteractionEventParticipant() {}

func (OrganizationEntity) ParticipantLabel() string {
	return NodeLabel_Organization
}

type OrganizationEntities []OrganizationEntity

func (organization OrganizationEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Organization,
		NodeLabel_Organization + "_" + tenant,
	}
}
