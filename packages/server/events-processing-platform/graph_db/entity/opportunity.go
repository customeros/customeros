package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type RenewalDetails struct {
	RenewedAt              *time.Time
	RenewalLikelihood      string
	RenewalUpdatedByUserId string
	RenewalUpdatedByUserAt *time.Time
}

type OpportunityEntity struct {
	Id                string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Source            neo4jentity.DataSource
	SourceOfTruth     neo4jentity.DataSource
	AppSource         string
	Name              string
	Amount            float64
	MaxAmount         float64
	InternalType      string
	ExternalType      string
	InternalStage     string
	ExternalStage     string
	EstimatedClosedAt *time.Time
	ClosedAt          *time.Time
	OwnerUserId       string
	CreatedByUserId   string
	GeneralNotes      string
	NextSteps         string
	Comments          string
	RenewalDetails    RenewalDetails
}
