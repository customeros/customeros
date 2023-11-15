package model

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

// Opportunity represents the state of an opportunity aggregate.
type Opportunity struct {
	ID                string                         `json:"id"`
	OrganizationId    string                         `json:"organizationId"`
	ContractId        string                         `json:"contractId"`
	Tenant            string                         `json:"tenant"`
	Name              string                         `json:"name"`
	Amount            float64                        `json:"amount"`
	InternalType      OpportunityInternalTypeString  `json:"internalType"`
	ExternalType      string                         `json:"externalType"`
	InternalStage     OpportunityInternalStageString `json:"internalStage"`
	ExternalStage     string                         `json:"externalStage"`
	EstimatedClosedAt *time.Time                     `json:"estimatedClosedAt,omitempty"`
	OwnerUserId       string                         `json:"ownerUserId"`
	CreatedByUserId   string                         `json:"createdByUserId"`
	Source            commonmodel.Source             `json:"source"`
	ExternalSystems   []commonmodel.ExternalSystem   `json:"externalSystems"`
	GeneralNotes      string                         `json:"generalNotes"`
	NextSteps         string                         `json:"nextSteps"`
	CreatedAt         time.Time                      `json:"createdAt"`
	UpdatedAt         time.Time                      `json:"updatedAt"`
}

// OpportunityDataFields contains all the fields that may be used to create or update an opportunity.
type OpportunityDataFields struct {
	Name              string
	Amount            float64
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
