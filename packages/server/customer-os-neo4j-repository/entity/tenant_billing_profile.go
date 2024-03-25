package entity

import (
	"time"
)

type TenantBillingProfileEntity struct {
	Id                     string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	LegalName              string
	Phone                  string
	AddressLine1           string
	AddressLine2           string
	AddressLine3           string
	Country                string
	Region                 string
	Locality               string
	Zip                    string
	VatNumber              string
	SendInvoicesFrom       string
	SendInvoicesBcc        string
	CanPayWithPigeon       bool
	CanPayWithBankTransfer bool
	Source                 DataSource
	SourceOfTruth          DataSource
	AppSource              string
	InvoiceNote            string
}

type TenantBillingProfileEntities []TenantBillingProfileEntity
