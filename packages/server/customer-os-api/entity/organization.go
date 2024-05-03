package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

// Deprecated, use neo4j module instead
type OrganizationEntity_TOBEDELETED struct {
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

func (OrganizationEntity_TOBEDELETED) IsNotedEntity() {}

func (OrganizationEntity_TOBEDELETED) NotedEntityLabel() string {
	return neo4jutil.NodeLabelOrganization
}

func (OrganizationEntity_TOBEDELETED) IsInteractionEventParticipant() {}

func (OrganizationEntity_TOBEDELETED) IsIssueParticipant() {}

func (OrganizationEntity_TOBEDELETED) ParticipantLabel() string {
	return neo4jutil.NodeLabelOrganization
}

func (OrganizationEntity_TOBEDELETED) IsMeetingParticipant() {}

func (OrganizationEntity_TOBEDELETED) MeetingParticipantLabel() string {
	return neo4jutil.NodeLabelOrganization
}

func (organization OrganizationEntity_TOBEDELETED) GetDataloaderKey() string {
	return organization.DataloaderKey
}

type OrganizationEntities_TOBEDELETED []OrganizationEntity_TOBEDELETED

func (organization OrganizationEntity_TOBEDELETED) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelOrganization,
		neo4jutil.NodeLabelOrganization + "_" + tenant,
	}
}
