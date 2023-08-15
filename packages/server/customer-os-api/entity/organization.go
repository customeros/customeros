package entity

import (
	"fmt"
	"time"
)

type RenewalLikelihoodProbability string

const (
	RenewalLikelihoodProbabilityHigh   RenewalLikelihoodProbability = "0-HIGH"
	RenewalLikelihoodProbabilityMedium RenewalLikelihoodProbability = "1-MEDIUM"
	RenewalLikelihoodProbabilityLow    RenewalLikelihoodProbability = "2-LOW"
	RenewalLikelihoodProbabilityZero   RenewalLikelihoodProbability = "3-ZERO"
)

type OrganizationEntity struct {
	ID                 string
	TenantOrganization bool
	Name               string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description        string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Website            string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry           string `neo4jDb:"property:industry;lookupName:INDUSTRY;supportCaseSensitive:true"`
	SubIndustry        string
	IndustryGroup      string
	TargetAudience     string
	ValueProposition   string
	IsPublic           bool
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	SlackChannelLink   string
	Employees          int64
	CreatedAt          time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	LastTouchpointAt   *time.Time `neo4jDb:"property:lastTouchpointAt;lookupName:LAST_TOUCHPOINT_AT;supportCaseSensitive:false"`
	LastTouchpointId   *string    `neo4jDb:"property:lastTouchpointId;lookupName:LAST_TOUCHPOINT_ID;supportCaseSensitive:false"`
	Source             DataSource
	SourceOfTruth      DataSource
	AppSource          string

	LinkedOrganizationType *string

	SuggestedMerge struct {
		SuggestedAt *time.Time
		SuggestedBy *string
		Confidence  *float64
	}
	RenewalLikelihood RenewalLikelihood
	RenewalForecast   RenewalForecast
	BillingDetails    BillingDetails

	InteractionEventParticipantDetails InteractionEventParticipantDetails

	DataloaderKey string
}

type RenewalLikelihood struct {
	RenewalLikelihood         string
	PreviousRenewalLikelihood string
	Comment                   *string
	UpdatedAt                 *time.Time
	UpdatedBy                 *string
}
type RenewalForecast struct {
	Amount          *float64
	PotentialAmount *float64
	Comment         *string
	UpdatedAt       *time.Time
	UpdatedBy       *string
}
type BillingDetails struct {
	Amount            *float64
	Frequency         string
	RenewalCycle      string
	RenewalCycleStart *time.Time
}

func (organization OrganizationEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organization.ID, organization.Name)
}

func (OrganizationEntity) IsNotedEntity() {}

func (OrganizationEntity) NotedEntityLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsInteractionEventParticipant() {}

func (OrganizationEntity) InteractionEventParticipantLabel() string {
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
