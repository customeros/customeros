package entity

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type InvoiceLineEntity struct {
	Id                      string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Name                    string
	Price                   float64
	Quantity                int64
	BilledType              neo4jenum.BilledType
	Amount                  float64
	Vat                     float64
	TotalAmount             float64
	Source                  DataSource
	SourceOfTruth           DataSource
	AppSource               string
	ServiceLineItemId       string
	ServiceLineItemParentId string

	DataloaderKey string
}

type InvoiceLineEntities []InvoiceLineEntity
