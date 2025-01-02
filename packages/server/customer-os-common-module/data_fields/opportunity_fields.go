package data_fields

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type OpportunityFields struct {
	AppSource         *string        `json:"appSource,omitempty"`
	Source            *string        `json:"source,omitempty"`
	CreatedAt         *time.Time     `json:"createdAt,omitempty"`
	Name              *string        `json:"name,omitempty"`
	Amount            *float64       `json:"amount,omitempty"`
	MaxAmount         *float64       `json:"maxAmount,omitempty"`
	ExternalStage     *string        `json:"externalStage,omitempty"`
	ExternalType      *string        `json:"externalType,omitempty"`
	EstimatedClosedAt *time.Time     `json:"estimatedClosedAt,omitempty"`
	InternalStage     *string        `json:"internalStage,omitempty"`
	InternalType      *string        `json:"internalType,omitempty"`
	Currency          *enum.Currency `json:"currency,omitempty"`
	NextSteps         *string        `json:"nextSteps,omitempty"`
	LikelihoodRate    *int64         `json:"likelihoodRate,omitempty"`
	OwnerId           *string        `json:"ownerId,omitempty"`
	OrganizationId    *string        `json:"organizationId,omitempty"`
	ContractId        *string        `json:"contractId,omitempty"`
}
