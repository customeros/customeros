package model

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

// Opportunity represents the state of an opportunity aggregate.
type Opportunity struct {
	ID                string                       `json:"id"`
	OrganizationId    string                       `json:"organizationId"`
	Tenant            string                       `json:"tenant"`
	Name              string                       `json:"name"`
	Amount            float64                      `json:"amount"`
	InternalType      OpportunityInternalType      `json:"internalType"`
	ExternalType      string                       `json:"externalType"`
	InternalStage     OpportunityInternalStage     `json:"internalStage"`
	ExternalStage     string                       `json:"externalStage"`
	EstimatedClosedAt *time.Time                   `json:"estimatedClosedAt,omitempty"`
	OwnerUserId       string                       `json:"ownerUserId"`
	CreatedByUserId   string                       `json:"createdByUserId"`
	Source            commonmodel.Source           `json:"source"`
	ExternalSystems   []commonmodel.ExternalSystem `json:"externalSystem"`
	GeneralNotes      string                       `json:"generalNotes"`
	NextSteps         string                       `json:"nextSteps"`
	CreatedAt         time.Time                    `json:"createdAt"`
	UpdatedAt         time.Time                    `json:"updatedAt"`
}

// OpportunityDataFields contains all the fields that may be used to create or update an opportunity.
type OpportunityDataFields struct {
	Name              string
	Amount            float64
	InternalType      OpportunityInternalTypeEnum
	ExternalType      string
	InternalStage     OpportunityInternalStageEnum
	ExternalStage     string
	EstimatedClosedAt *time.Time
	OwnerUserId       string
	CreatedByUserId   string
	GeneralNotes      string
	NextSteps         string
	OrganizationId    string `json:"organizationId" validate:"required"`
}

// OpportunityInternalTypeEnum represents the type of opportunity within the system.
type OpportunityInternalTypeEnum int32

const (
	NBO OpportunityInternalTypeEnum = iota
	UPSELL
	CROSS_SELL
)

type OpportunityInternalType string

const (
	OpportunityInternalTypeNBO       OpportunityInternalType = "NBO"
	OpportunityInternalTypeUpsell                            = "UPSELL"
	OpportunityInternalTypeCrossSell                         = "CROSS_SELL"
)

// String returns the string representation of the OpportunityInternalTypeEnum.
func (t OpportunityInternalTypeEnum) StrValue() OpportunityInternalType {
	switch t {
	case NBO:
		return OpportunityInternalTypeNBO
	case UPSELL:
		return OpportunityInternalTypeUpsell
	case CROSS_SELL:
		return OpportunityInternalTypeCrossSell
	default:
		return ""
	}
}

// OpportunityInternalStageEnum represents the stage of the opportunity within the system.
type OpportunityInternalStageEnum int32

const (
	OPEN OpportunityInternalStageEnum = iota
	EVALUATING
	CLOSED_WON
	CLOSED_LOST
)

type OpportunityInternalStage string

const (
	OpportunityInternalStageOpen       OpportunityInternalStage = "OPEN"
	OpportunityInternalStageEvaluating                          = "EVALUATING"
	OpportunityInternalStageClosedWon                           = "CLOSED_WON"
	OpportunityInternalStageClosedLost                          = "CLOSED_LOST"
)

// String returns the string representation of the OpportunityInternalStageEnum.
func (s OpportunityInternalStageEnum) StrValue() OpportunityInternalStage {
	switch s {
	case OPEN:
		return OpportunityInternalStageOpen
	case EVALUATING:
		return OpportunityInternalStageEvaluating
	case CLOSED_WON:
		return OpportunityInternalStageClosedWon
	case CLOSED_LOST:
		return OpportunityInternalStageClosedLost
	default:
		return ""
	}
}
