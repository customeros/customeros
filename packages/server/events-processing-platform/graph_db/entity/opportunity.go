package entity

import (
	"time"
)

type RenewalDetails struct {
	RenewedAt         *time.Time
	RenewalLikelihood string
}

type OpportunityEntity struct {
	Id                string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Source            DataSource
	SourceOfTruth     DataSource
	AppSource         string
	Name              string
	Amount            float64
	MaxAmount         float64
	InternalType      string
	ExternalType      string
	InternalStage     string
	ExternalStage     string
	EstimatedClosedAt *time.Time
	OwnerUserId       string
	CreatedByUserId   string
	GeneralNotes      string
	NextSteps         string
	RenewalDetails    RenewalDetails
}
