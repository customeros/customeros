package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type InvoiceEntity struct {
	Id                            string
	CreatedAt                     time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt                     time.Time
	DryRun                        bool          `neo4jDb:"property:dryRun;lookupName:DRY_RUN;supportCaseSensitive:false"`
	Number                        string        `neo4jDb:"property:number;lookupName:NUMBER;supportCaseSensitive:false"`
	Currency                      enum.Currency `neo4jDb:"property:currency;lookupName:CURRENCY;supportCaseSensitive:false"`
	PeriodStartDate               time.Time
	PeriodEndDate                 time.Time
	DueDate                       time.Time
	DomesticPaymentsBankInfo      string
	InternationalPaymentsBankInfo string
	Customer                      InvoiceCustomer
	Provider                      InvoiceProvider
	Amount                        float64 `neo4jDb:"property:amount;lookupName:AMOUNT;supportCaseSensitive:false"`
	Vat                           float64 `neo4jDb:"property:vat;lookupName:VAT;supportCaseSensitive:false"`
	SubtotalAmount                float64 `neo4jDb:"property:subtotalAmount;lookupName:SUBTOTAL_AMOUNT;supportCaseSensitive:false"`
	TotalAmount                   float64 `neo4jDb:"property:totalAmount;lookupName:TOTAL_AMOUNT;supportCaseSensitive:false"`
	RepositoryFileId              string
	BillingCycle                  enum.BillingCycle
	Status                        enum.InvoiceStatus
	Note                          string
	PaymentDetails                PaymentDetails

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	InvoiceInternalFields InvoiceInternalFields
}

type InvoiceCustomer struct {
	Name         string
	Email        string
	AddressLine1 string
	AddressLine2 string
	Zip          string
	Locality     string
	Country      string
}

type InvoiceProvider struct {
	LogoUrl      string
	Name         string
	Email        string
	AddressLine1 string
	AddressLine2 string
	Zip          string
	Locality     string
	Country      string
}

type PaymentDetails struct {
	PaymentLink string
}

type InvoiceInternalFields struct {
	PaymentRequestedAt                *time.Time
	PayInvoiceNotificationRequestedAt *time.Time // used for locking in batch to not send the same notification multiple times under an hour
	PayInvoiceNotificationSentAt      *time.Time // used to prevent sending the same notification
	PaidInvoiceNotificationSentAt     *time.Time // not used, but set for now
}

type InvoiceEntities []InvoiceEntity
