package entity

import (
	"time"
)

type TenantBillingProfile struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LegalName     string
	Email         string
	Phone         string
	AddressLine1  string
	AddressLine2  string
	AddressLine3  string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}
