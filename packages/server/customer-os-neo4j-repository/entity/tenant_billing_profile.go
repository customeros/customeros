package entity

import (
	"time"
)

type TenantBillingProfileEntity struct {
	Id                                string
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
	LegalName                         string
	Phone                             string
	AddressLine1                      string
	AddressLine2                      string
	AddressLine3                      string
	Locality                          string
	Country                           string
	Zip                               string
	DomesticPaymentsBankInfo          string
	InternationalPaymentsBankInfo     string
	DomesticPaymentsBankName          string
	DomesticPaymentsAccountNumber     string
	DomesticPaymentsSortCode          string
	InternationalPaymentsSwiftBic     string
	InternationalPaymentsBankName     string
	InternationalPaymentsBankAddress  string
	InternationalPaymentsInstructions string
	VatNumber                         string
	SendInvoicesFrom                  string
	SendInvoicesBcc                   string
	CanPayWithCard                    bool
	CanPayWithDirectDebitSEPA         bool
	CanPayWithDirectDebitACH          bool
	CanPayWithDirectDebitBacs         bool
	CanPayWithPigeon                  bool
	Source                            DataSource
	SourceOfTruth                     DataSource
	AppSource                         string
}

type TenantBillingProfileEntities []TenantBillingProfileEntity
