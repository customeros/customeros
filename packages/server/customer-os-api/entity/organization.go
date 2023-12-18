package entity

import (
	"fmt"
	"time"
)

type OrganizationEntity struct {
	ID                 string
	CustomerOsId       string `neo4jDb:"property:customerOsId;lookupName:CUSTOMER_OS_ID;supportCaseSensitive:false"`
	ReferenceId        string `neo4jDb:"property:referenceId;lookupName:REFERENCE_ID;supportCaseSensitive:true"`
	Name               string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description        string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Website            string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry           string `neo4jDb:"property:industry;lookupName:INDUSTRY;supportCaseSensitive:true"`
	SubIndustry        string
	IndustryGroup      string
	TargetAudience     string
	ValueProposition   string
	IsPublic           bool
	IsCustomer         bool `neo4jDb:"property:isCustomer;lookupName:IS_CUSTOMER;supportCaseSensitive:false"`
	Hide               bool
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	Note               string
	Employees          int64
	CreatedAt          time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	LastTouchpointId   *string    `neo4jDb:"property:lastTouchpointId;lookupName:LAST_TOUCHPOINT_ID;supportCaseSensitive:false"`
	LastTouchpointAt   *time.Time `neo4jDb:"property:lastTouchpointAt;lookupName:LAST_TOUCHPOINT_AT;supportCaseSensitive:false"`
	LastTouchpointType *string    `neo4jDb:"property:lastTouchpointType;lookupName:LAST_TOUCHPOINT_TYPE;supportCaseSensitive:false"`
	YearFounded        *int64
	Headquarters       string
	EmployeeGrowthRate string
	LogoUrl            string
	Source             DataSource
	SourceOfTruth      DataSource
	AppSource          string

	LinkedOrganizationType *string

	SuggestedMerge struct {
		SuggestedAt *time.Time
		SuggestedBy *string
		Confidence  *float64
	}
	RenewalSummary    RenewalSummary
	OnboardingDetails OnboardingDetails

	InteractionEventParticipantDetails InteractionEventParticipantDetails

	DataloaderKey string
}

type RenewalSummary struct {
	ArrForecast            *float64
	MaxArrForecast         *float64
	NextRenewalAt          *time.Time `neo4jDb:"property:derivedNextRenewalAt;lookupName:RENEWAL_DATE;supportCaseSensitive:false"`
	RenewalLikelihood      string     `neo4jDb:"property:derivedRenewalLikelihood;lookupName:RENEWAL_LIKELIHOOD;supportCaseSensitive:false"`
	RenewalLikelihoodOrder *int64
}

type OnboardingDetails struct {
	Status       OnboardingStatus
	SortingOrder *int64
	UpdatedAt    *time.Time
	Comments     string
}

func (organization OrganizationEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organization.ID, organization.Name)
}

func (OrganizationEntity) IsNotedEntity() {}

func (OrganizationEntity) NotedEntityLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsInteractionEventParticipant() {}

func (OrganizationEntity) IsIssueParticipant() {}

func (OrganizationEntity) ParticipantLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsMeetingParticipant() {}

func (OrganizationEntity) MeetingParticipantLabel() string {
	return NodeLabel_Organization
}

func (organization OrganizationEntity) GetDataloaderKey() string {
	return organization.DataloaderKey
}

type OrganizationEntities []OrganizationEntity

func (organization OrganizationEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Organization,
		NodeLabel_Organization + "_" + tenant,
	}
}
