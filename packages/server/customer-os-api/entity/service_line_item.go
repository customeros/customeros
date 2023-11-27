package entity

import "time"

type ServiceLineItemEntity struct {
	ID            string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	EndedAt       *time.Time
	Billed        BilledType
	Price         float64
	Quantity      int64
	Comments      string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

type ServiceLineItemEntities []ServiceLineItemEntity
