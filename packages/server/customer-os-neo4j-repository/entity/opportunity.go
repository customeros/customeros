package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type OpportunityProperty string

const (
	OpportunityPropertyAmount         OpportunityProperty = "amount"
	OpportunityPropertyMaxAmount      OpportunityProperty = "maxAmount"
	OpportunityPropertyNextSteps      OpportunityProperty = "nextSteps"
	OpportunityPropertyCurrency       OpportunityProperty = "currency"
	OpportunityPropertyLikelihoodRate OpportunityProperty = "likelihoodRate"
	OpportunityPropertyStageUpdatedAt OpportunityProperty = "stageUpdatedAt"
)

type OpportunityEntity struct {
	DataLoaderKey
	Id                string
	Name              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Amount            float64
	MaxAmount         float64
	InternalType      enum.OpportunityInternalType
	ExternalType      string
	InternalStage     enum.OpportunityInternalStage
	ExternalStage     string
	EstimatedClosedAt *time.Time
	ClosedAt          *time.Time
	GeneralNotes      string
	NextSteps         string
	Comments          string
	Source            DataSource
	SourceOfTruth     DataSource
	AppSource         string
	OwnerUserId       string
	RenewalDetails    RenewalDetails
	InternalFields    OpportunityInternalFields
	LikelihoodRate    int64
	Currency          enum.Currency
	StageUpdatedAt    *time.Time
}

type OpportunityInternalFields struct {
	RolloutRenewalRequestedAt *time.Time
}

type RenewalDetails struct {
	RenewedAt              *time.Time // DateTime
	RenewalLikelihood      enum.RenewalLikelihood
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt *time.Time
	RenewalApproved        bool
	RenewalAdjustedRate    int64
}

type OpportunityEntities []OpportunityEntity
