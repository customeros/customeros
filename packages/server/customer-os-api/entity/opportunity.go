package entity

import (
	"time"
)

type OpportunityEntity struct {
	Id                     string
	Name                   string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Amount                 float64
	MaxAmount              float64
	InternalType           InternalType
	ExternalType           string
	InternalStage          InternalStage
	ExternalStage          string
	EstimatedClosedAt      time.Time
	GeneralNotes           string
	NextSteps              string
	RenewedAt              time.Time
	RenewalLikelihood      OpportunityRenewalLikelihood
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt time.Time
	Comments               string
	Source                 DataSource
	SourceOfTruth          DataSource
	AppSource              string

	DataloaderKey string
}

type OpportunityRenewalEntity struct {
	Id                     string
	Name                   string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Amount                 float64
	RenewalLikelihood      OpportunityRenewalLikelihood
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt time.Time
	Comments               string
	Source                 DataSource
	SourceOfTruth          DataSource
	AppSource              string

	DataloaderKey string
}

type OpportunityEntities []OpportunityEntity
