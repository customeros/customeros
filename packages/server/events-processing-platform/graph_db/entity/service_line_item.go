package entity

import (
	"time"
)

type ServiceLineItemEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Name          string
	Billed        string
	Price         float64
	Quantity      int64
	Comments      string
}

type ServiceLineItemEntities []ServiceLineItemEntity
