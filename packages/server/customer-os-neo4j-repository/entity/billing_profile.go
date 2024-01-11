package entity

import (
	"time"
)

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
