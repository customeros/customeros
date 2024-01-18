package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type InvoiceEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	DryRun           bool
	Number           string
	Currency         enum.Currency
	Date             time.Time
	DueDate          time.Time
	Amount           float64
	Vat              float64
	Total            float64
	RepositoryFileId string

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	InvoiceInternalFields InvoiceInternalFields
}

type InvoiceInternalFields struct {
	PaymentRequestedAt *time.Time
}

type InvoiceEntities []InvoiceEntity
