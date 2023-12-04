package entity

import (
	"fmt"
	"time"
)

type OrganizationEntity struct {
	ID                string
	CustomerOsId      string
	Name              string
	Description       string
	Website           string
	Industry          string
	SubIndustry       string
	IndustryGroup     string
	TargetAudience    string
	ValueProposition  string
	IsPublic          bool
	IsCustomer        bool
	Hide              bool
	Market            string
	LastFundingRound  string
	LastFundingAmount string
	ReferenceId       string
	Note              string
	Employees         int64
	CreatedAt         time.Time
	LastTouchpointAt  *time.Time
	UpdatedAt         time.Time
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
	RenewalSummary RenewalSummary

	InteractionEventParticipantDetails InteractionEventParticipantDetails

	DataloaderKey string
}

type RenewalSummary struct {
	ArrForecast            *float64
	MaxArrForecast         *float64
	NextRenewalAt          *time.Time
	RenewalLikelihood      string
	RenewalLikelihoodOrder *int64
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
