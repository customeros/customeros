package entity

import (
	"time"
)

type BillingProfileEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type BillingProfileEntities []BillingProfileEntity
