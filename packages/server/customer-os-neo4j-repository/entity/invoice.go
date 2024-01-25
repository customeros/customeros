package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type InvoiceEntity struct {
	Id               string
	CreatedAt        time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt        time.Time
	DryRun           bool
	Number           string        `neo4jDb:"property:number;lookupName:NUMBER;supportCaseSensitive:false"`
	Currency         enum.Currency `neo4jDb:"property:currency;lookupName:CURRENCY;supportCaseSensitive:false"`
	PeriodStartDate  time.Time
	PeriodEndDate    time.Time
	DueDate          time.Time
	Amount           float64 `neo4jDb:"property:amount;lookupName:AMOUNT;supportCaseSensitive:false"`
	Vat              float64 `neo4jDb:"property:vat;lookupName:VAT;supportCaseSensitive:false"`
	TotalAmount      float64 `neo4jDb:"property:totalAmount;lookupName:TOTAL_AMOUNT;supportCaseSensitive:false"`
	RepositoryFileId string
	BillingCycle     enum.BillingCycle
	Status           enum.InvoiceStatus
	Note             string

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	InvoiceInternalFields InvoiceInternalFields
}

type InvoiceInternalFields struct {
	PaymentRequestedAt *time.Time
}

type InvoiceEntities []InvoiceEntity
