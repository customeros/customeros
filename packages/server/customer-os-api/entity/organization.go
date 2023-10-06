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
	ID                string
	CustomerOsId      string `neo4jDb:"property:customerOsId;lookupName:CUSTOMER_OS_ID;supportCaseSensitive:false"`
	ReferenceId       string `neo4jDb:"property:referenceId;lookupName:REFERENCE_ID;supportCaseSensitive:true"`
	Name              string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description       string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Website           string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry          string `neo4jDb:"property:industry;lookupName:INDUSTRY;supportCaseSensitive:true"`
	SubIndustry       string
	IndustryGroup     string
	TargetAudience    string
	ValueProposition  string
	IsPublic          bool
	IsCustomer        bool `neo4jDb:"property:isCustomer;lookupName:IS_CUSTOMER;supportCaseSensitive:false"`
	Hide              bool
	Market            string
	LastFundingRound  string
	LastFundingAmount string
	Note              string
	Employees         int64
	CreatedAt         time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt         time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	LastTouchpointAt  *time.Time `neo4jDb:"property:lastTouchpointAt;lookupName:LAST_TOUCHPOINT_AT;supportCaseSensitive:false"`
	LastTouchpointId  *string    `neo4jDb:"property:lastTouchpointId;lookupName:LAST_TOUCHPOINT_ID;supportCaseSensitive:false"`
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
	RenewalLikelihood         string `neo4jDb:"property:renewalLikelihood;lookupName:RENEWAL_LIKELIHOOD;supportCaseSensitive:false"`
	PreviousRenewalLikelihood string
	Comment                   *string
	UpdatedAt                 *time.Time
	UpdatedBy                 *string
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
	if r.UpdatedBy != nil {
		output += fmt.Sprintf(", UpdatedBy: %v", *r.UpdatedBy)
	} else {
		output += ", UpdatedBy: nil"
	}
	return output
}

type RenewalForecast struct {
	Amount          *float64 `neo4jDb:"property:renewalForecastAmount;lookupName:FORECAST_AMOUNT;supportCaseSensitive:false"`
	PotentialAmount *float64
	Comment         *string
	UpdatedAt       *time.Time
	UpdatedById     *string
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
	if r.UpdatedById != nil {
		output += fmt.Sprintf(", UpdatedById: %v", *r.UpdatedById)
	} else {
		output += ", UpdatedById: nil"
	}
	return output
}

type BillingDetails struct {
	Amount            *float64
	Frequency         string
	RenewalCycle      string
	RenewalCycleStart *time.Time
	RenewalCycleNext  *time.Time `neo4jDb:"property:billingDetailsRenewalCycleNext;lookupName:RENEWAL_CYCLE_NEXT;supportCaseSensitive:false"`
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
