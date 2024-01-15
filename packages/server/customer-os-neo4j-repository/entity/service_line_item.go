package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ServiceLineItemEntity struct {
	ID               string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	StartedAt        time.Time
	EndedAt          *time.Time
	IsCanceled       bool
	Billed           enum.BilledType
	Price            float64
	Quantity         int64
	PreviousBilled   enum.BilledType
	PreviousPrice    float64
	PreviousQuantity int64
	Comments         string
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	ParentID         string
}
