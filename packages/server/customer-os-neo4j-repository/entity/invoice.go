package entity

import (
	"time"
)

type InvoiceEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	DryRun           bool
	Number           string
	Date             time.Time
	DueDate          time.Time
	Amount           float64
	Vat              float64
	Total            float64
	RepositoryFileId string
	PdfGenerated     bool

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}
