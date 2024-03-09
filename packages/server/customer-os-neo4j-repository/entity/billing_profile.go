package entity

import (
	"time"
)

// Deprecated - to be checked of not used and remove it
type BillingProfileEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LegalName     string
	TaxId         string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type BillingProfileEntities []BillingProfileEntity
