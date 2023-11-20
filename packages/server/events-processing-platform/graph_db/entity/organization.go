package entity

import (
	"fmt"
	"time"
)

type RenewalLikelihoodProbability string

const (
	RenewalLikelihoodHigh   RenewalLikelihoodProbability = "0-HIGH"
	RenewalLikelihoodMedium RenewalLikelihoodProbability = "1-MEDIUM"
	RenewalLikelihoodLow    RenewalLikelihoodProbability = "2-LOW"
	RenewalLikelihoodZero   RenewalLikelihoodProbability = "3-ZERO"
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
	UpdatedBy                 string
}

func (r RenewalLikelihood) String() string {
	output := ""
	output += fmt.Sprintf("RenewalLikelihood: %v, Previous: %v", r.RenewalLikelihood, r.PreviousRenewalLikelihood)
	if r.Comment != nil {
		output += fmt.Sprintf(", Comment: %v", *r.Comment)
	} else {
		output += ", Comment: nil"
	}
	if r.UpdatedAt != nil {
		output += fmt.Sprintf(", UpdatedAt: %v", *r.UpdatedAt)
	} else {
		output += ", UpdatedAt: nil"
	}
	output += fmt.Sprintf(", UpdatedBy: %v", r.UpdatedBy)
	return output
}

type RenewalForecast struct {
	Amount          *float64
	PotentialAmount *float64
	Comment         *string
	UpdatedAt       *time.Time
	UpdatedBy       string
	Arr             *float64
	MaxArr          *float64
}

func (r RenewalForecast) String() string {
	output := ""
	if r.Amount != nil {
		output += fmt.Sprintf("Amount: %v", *r.Amount)
	} else {
		output += "Amount: nil"
	}
	if r.PotentialAmount != nil {
		output += fmt.Sprintf(", Potential: %v", *r.PotentialAmount)
	} else {
		output += ", Potential: nil"
	}
	if r.Comment != nil {
		output += fmt.Sprintf(", Comment: %v", *r.Comment)
	} else {
		output += ", Comment: nil"
	}
	if r.UpdatedAt != nil {
		output += fmt.Sprintf(", UpdatedAt: %v", *r.UpdatedAt)
	} else {
		output += ", UpdatedAt: nil"
	}
	output += fmt.Sprintf(", UpdatedBy: %v", r.UpdatedBy)
	return output
}

type BillingDetails struct {
	Amount            *float64
	Frequency         string
	RenewalCycle      string
	RenewalCycleStart *time.Time
	RenewalCycleNext  *time.Time
}

func (b BillingDetails) String() string {
	output := ""
	if b.Amount != nil {
		output += fmt.Sprintf("Amount: %v", *b.Amount)
	} else {
		output += "Amount: nil"
	}
	output += fmt.Sprintf(", Frequency: %v", b.Frequency)
	output += fmt.Sprintf(", RenewalCycle: %v", b.RenewalCycle)
	if b.RenewalCycleStart != nil {
		output += fmt.Sprintf(", RenewalCycleStart: %v", *b.RenewalCycleStart)
	} else {
		output += ", RenewalCycleStart: nil"
	}
	if b.RenewalCycleNext != nil {
		output += fmt.Sprintf(", RenewalCycleNext: %v", *b.RenewalCycleNext)
	} else {
		output += ", RenewalCycleNext: nil"
	}
	return output
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
