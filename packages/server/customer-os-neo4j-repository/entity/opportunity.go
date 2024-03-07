package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type OpportunityEntity struct {
	Id                     string
	Name                   string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Amount                 float64
	MaxAmount              float64
	InternalType           enum.InternalType
	ExternalType           string
	InternalStage          enum.InternalStage
	ExternalStage          string
	EstimatedClosedAt      *time.Time
	ClosedAt               *time.Time
	GeneralNotes           string
	NextSteps              string
	RenewedAt              time.Time
	RenewalLikelihood      enum.RenewalLikelihood
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt time.Time
	Comments               string
	Source                 DataSource
	SourceOfTruth          DataSource
	AppSource              string
	OwnerUserId            string
}
