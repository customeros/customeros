package model

import (
	"time"

	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

const (
	FieldMaskName              = "name"
	FieldMaskAmount            = "amount"
	FieldMaskMaxAmount         = "maxAmount"
	FieldMaskComments          = "comments"
	FieldMaskRenewalLikelihood = "renewalLikelihood"
	FieldMaskRenewalApproved   = "renewalApproved"
	FieldMaskRenewedAt         = "renewedAt"
	FieldMaskAdjustedRate      = "adjustedRate"
	FieldMaskExternalType      = "externalType"
	FieldMaskInternalType      = "internalType"
	FieldMaskExternalStage     = "externalStage"
	FieldMaskInternalStage     = "internalStage"
	FieldMaskEstimatedClosedAt = "estimatedClosedAt"
)

type RenewalDetails struct {
	RenewedAt              *time.Time `json:"renewedAt,omitempty"`
	RenewalLikelihood      string     `json:"renewalLikelihood,omitempty"`
	RenewalUpdatedByUserId string     `json:"renewalUpdatedByUserId,omitempty"`
	RenewalUpdatedByUserAt *time.Time `json:"renewalUpdatedByUserAt,omitempty"`
	RenewalApproved        bool       `json:"renewalApproved,omitempty"`
	RenewalAdjustedRate    int64      `json:"renewalAdjustedRate,omitempty"`
}

// Opportunity represents the state of an opportunity aggregate.
type Opportunity struct {
	ID                string                       `json:"id"`
	OrganizationId    string                       `json:"organizationId"`
	ContractId        string                       `json:"contractId"`
	Tenant            string                       `json:"tenant"`
	Name              string                       `json:"name"`
	Amount            float64                      `json:"amount"`
	MaxAmount         float64                      `json:"maxAmount"`
	InternalType      string                       `json:"internalType"`
	ExternalType      string                       `json:"externalType"`
	InternalStage     string                       `json:"internalStage"`
	ExternalStage     string                       `json:"externalStage"`
	EstimatedClosedAt *time.Time                   `json:"estimatedClosedAt,omitempty"`
	ClosedAt          *time.Time                   `json:"closedAt,omitempty"`
	OwnerUserId       string                       `json:"ownerUserId"`
	CreatedByUserId   string                       `json:"createdByUserId"`
	Source            commonmodel.Source           `json:"source"`
	ExternalSystems   []commonmodel.ExternalSystem `json:"externalSystems"`
	GeneralNotes      string                       `json:"generalNotes"`
	NextSteps         string                       `json:"nextSteps"`
	CreatedAt         time.Time                    `json:"createdAt"`
	UpdatedAt         time.Time                    `json:"updatedAt"`
	RenewalDetails    RenewalDetails               `json:"renewal,omitempty"`
	Comments          string                       `json:"comments,omitempty"`
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
	RenewedAt         *time.Time
}

// OpportunityInternalType represents the type of opportunity within the system.
type OpportunityInternalType int32

const (
	NBO OpportunityInternalType = iota
	UPSELL
	CROSS_SELL
	RENEWAL
)

// String returns the string representation of the OpportunityInternalType.
func (t OpportunityInternalType) StringEnumValue() neo4jenum.OpportunityInternalType {
	switch t {
	case NBO:
		return neo4jenum.OpportunityInternalTypeNBO
	case UPSELL:
		return neo4jenum.OpportunityInternalTypeUpsell
	case CROSS_SELL:
		return neo4jenum.OpportunityInternalTypeCrossSell
	case RENEWAL:
		return neo4jenum.OpportunityInternalTypeRenewal
	default:
		return ""
	}
}

// OpportunityInternalStage represents the stage of the opportunity within the system.
type OpportunityInternalStage int32

const (
	OPEN OpportunityInternalStage = iota
	CLOSED_WON
	CLOSED_LOST
)

func (t OpportunityInternalStage) StringEnumValue() neo4jenum.OpportunityInternalStage {
	switch t {
	case OPEN:
		return neo4jenum.OpportunityInternalStageOpen
	case CLOSED_WON:
		return neo4jenum.OpportunityInternalStageClosedWon
	case CLOSED_LOST:
		return neo4jenum.OpportunityInternalStageClosedLost
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
func (r RenewalLikelihood) StringEnumValue() neo4jenum.RenewalLikelihood {
	switch r {
	case HIGH_RENEWAL:
		return neo4jenum.RenewalLikelihoodHigh
	case MEDIUM_RENEWAL:
		return neo4jenum.RenewalLikelihoodMedium
	case LOW_RENEWAL:
		return neo4jenum.RenewalLikelihoodLow
	case ZERO_RENEWAL:
		return neo4jenum.RenewalLikelihoodZero
	default:
		return ""
	}
}
