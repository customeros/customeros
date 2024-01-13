package entity

import (
	"time"
)

type InvoiceLineEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Index    int64
	Name     string
	Price    float64
	Quantity int64
	Amount   float64
	Vat      float64
	Total    float64

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}
