package entity

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
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
	EstimatedClosedAt      *time.Time
	GeneralNotes           string
	NextSteps              string
	RenewedAt              time.Time
	RenewalLikelihood      OpportunityRenewalLikelihood
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt time.Time
	Comments               string
	Source                 neo4jentity.DataSource
	SourceOfTruth          neo4jentity.DataSource
	AppSource              string
	OwnerUserId            string

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
	Source                 neo4jentity.DataSource
	SourceOfTruth          neo4jentity.DataSource
	AppSource              string

	DataloaderKey string
}

type OpportunityEntities []OpportunityEntity
