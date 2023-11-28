package model

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

const (
	FieldMaskName      = "name"
	FieldMaskAmount    = "amount"
	FieldMaskMaxAmount = "maxAmount"
)

type RenewalDetails struct {
	RenewedAt              *time.Time `json:"renewedAt,omitempty"`
	RenewalLikelihood      string     `json:"renewalLikelihood,omitempty"`
	RenewalUpdatedByUserId string     `json:"renewalUpdatedByUserId,omitempty"`
	RenewalUpdatedByUserAt *time.Time `json:"renewalUpdatedByUserAt,omitempty"`
}

// Opportunity represents the state of an opportunity aggregate.
type Opportunity struct {
	ID                string                         `json:"id"`
	OrganizationId    string                         `json:"organizationId"`
	ContractId        string                         `json:"contractId"`
	Tenant            string                         `json:"tenant"`
	Name              string                         `json:"name"`
	Amount            float64                        `json:"amount"`
	MaxAmount         float64                        `json:"maxAmount"`
	InternalType      OpportunityInternalTypeString  `json:"internalType"`
	ExternalType      string                         `json:"externalType"`
	InternalStage     OpportunityInternalStageString `json:"internalStage"`
	ExternalStage     string                         `json:"externalStage"`
	EstimatedClosedAt *time.Time                     `json:"estimatedClosedAt,omitempty"`
	ClosedAt          *time.Time                     `json:"closedAt,omitempty"`
	OwnerUserId       string                         `json:"ownerUserId"`
	CreatedByUserId   string                         `json:"createdByUserId"`
	Source            commonmodel.Source             `json:"source"`
	ExternalSystems   []commonmodel.ExternalSystem   `json:"externalSystems"`
	GeneralNotes      string                         `json:"generalNotes"`
	NextSteps         string                         `json:"nextSteps"`
	CreatedAt         time.Time                      `json:"createdAt"`
	UpdatedAt         time.Time                      `json:"updatedAt"`
	RenewalDetails    RenewalDetails                 `json:"renewal,omitempty"`
	Comments          string                         `json:"comments,omitempty"`
}

// OpportunityDataFields contains all the fields that may be used to create or update an opportunity.
type OpportunityDataFields struct {
	Name              string
	Amount            float64
	MaxAmount         float64
	InternalType      OpportunityInternalType
	ExternalType      string
	InternalStage     OpportunityInternalStage
	ExternalStage     string
	EstimatedClosedAt *time.Time
	OwnerUserId       string
	CreatedByUserId   string
	GeneralNotes      string
	NextSteps         string
	OrganizationId    string
}

// OpportunityInternalType represents the type of opportunity within the system.
type OpportunityInternalType int32

const (
	NBO OpportunityInternalType = iota
	UPSELL
	CROSS_SELL
	RENEWAL
)

type OpportunityInternalTypeString string

const (
	OpportunityInternalTypeStringNBO       OpportunityInternalTypeString = "NBO"
	OpportunityInternalTypeStringUpsell    OpportunityInternalTypeString = "UPSELL"
	OpportunityInternalTypeStringCrossSell OpportunityInternalTypeString = "CROSS_SELL"
	OpportunityInternalTypeStringRenewal   OpportunityInternalTypeString = "RENEWAL"
)

// String returns the string representation of the OpportunityInternalType.
func (t OpportunityInternalType) StringValue() OpportunityInternalTypeString {
	switch t {
	case NBO:
		return OpportunityInternalTypeStringNBO
	case UPSELL:
		return OpportunityInternalTypeStringUpsell
	case CROSS_SELL:
		return OpportunityInternalTypeStringCrossSell
	case RENEWAL:
		return OpportunityInternalTypeStringRenewal
	default:
		return ""
	}
}

func OpportunityInternalTypeStringDecode(val string) OpportunityInternalTypeString {
	switch val {
	case "NBO":
		return OpportunityInternalTypeStringNBO
	case "UPSELL":
		return OpportunityInternalTypeStringUpsell
	case "CROSS_SELL":
		return OpportunityInternalTypeStringCrossSell
	case "RENEWAL":
		return OpportunityInternalTypeStringRenewal
	default:
		return ""
	}
}

// OpportunityInternalStage represents the stage of the opportunity within the system.
type OpportunityInternalStage int32

const (
	OPEN OpportunityInternalStage = iota
	EVALUATING
	CLOSED_WON
	CLOSED_LOST
)

type OpportunityInternalStageString string

const (
	OpportunityInternalStageStringOpen       OpportunityInternalStageString = "OPEN"
	OpportunityInternalStageStringEvaluating OpportunityInternalStageString = "EVALUATING"
	OpportunityInternalStageStringClosedWon  OpportunityInternalStageString = "CLOSED_WON"
	OpportunityInternalStageStringClosedLost OpportunityInternalStageString = "CLOSED_LOST"
)

// String returns the string representation of the OpportunityInternalStage.
func (s OpportunityInternalStage) StringValue() OpportunityInternalStageString {
	switch s {
	case OPEN:
		return OpportunityInternalStageStringOpen
	case EVALUATING:
		return OpportunityInternalStageStringEvaluating
	case CLOSED_WON:
		return OpportunityInternalStageStringClosedWon
	case CLOSED_LOST:
		return OpportunityInternalStageStringClosedLost
	default:
		return ""
	}
}

func OpportunityInternalStageStringDecode(val string) OpportunityInternalStageString {
	switch val {
	case "OPEN":
		return OpportunityInternalStageStringOpen
	case "EVALUATING":
		return OpportunityInternalStageStringEvaluating
	case "CLOSED_WON":
		return OpportunityInternalStageStringClosedWon
	case "CLOSED_LOST":
		return OpportunityInternalStageStringClosedLost
	default:
		return ""
	}
}

// OpportunityInternalType represents the type of opportunity within the system.
type RenewalLikelihood int32

const (
	HIGH_RENEWAL RenewalLikelihood = iota
	MEDIUM_RENEWAL
	LOW_RENEWAL
	ZERO_RENEWAL
)

// String returns the string representation of the OpportunityInternalStage.
func (r RenewalLikelihood) StringValue() RenewalLikelihoodString {
	switch r {
	case HIGH_RENEWAL:
		return RenewalLikelihoodStringHigh
	case MEDIUM_RENEWAL:
		return RenewalLikelihoodStringMedium
	case LOW_RENEWAL:
		return RenewalLikelihoodStringLow
	case ZERO_RENEWAL:
		return RenewalLikelihoodStringZero
	default:
		return ""
	}
}

type RenewalLikelihoodString string

const (
	RenewalLikelihoodStringHigh   RenewalLikelihoodString = "HIGH_RENEWAL"
	RenewalLikelihoodStringMedium RenewalLikelihoodString = "MEDIUM_RENEWAL"
	RenewalLikelihoodStringLow    RenewalLikelihoodString = "LOW_RENEWAL"
	RenewalLikelihoodStringZero   RenewalLikelihoodString = "ZERO_RENEWAL"
)

func RenewalLikelihoodStringDecode(val string) RenewalLikelihoodString {
	switch val {
	case "HIGH_RENEWAL":
		return RenewalLikelihoodStringHigh
	case "MEDIUM_RENEWAL":
		return RenewalLikelihoodStringMedium
	case "LOW_RENEWAL":
		return RenewalLikelihoodStringLow
	case "ZERO_RENEWAL":
		return RenewalLikelihoodStringZero
	default:
		return ""
	}
}
