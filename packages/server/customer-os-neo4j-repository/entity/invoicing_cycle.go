package entity

import (
	"time"
)

type InvoicingCycleEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Type          InvoicingCycleType
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type InvoicingCycleType string

const (
	InvoicingCycleTypeDate        InvoicingCycleType = "DATE"
	InvoicingCycleTypeAnniversary InvoicingCycleType = "ANNIVERSARY"
)
