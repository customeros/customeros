package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
	Billed           BilledType
	Price            float64
	Quantity         int64
	PreviousBilled   BilledType
	PreviousPrice    float64
	PreviousQuantity int64
	Comments         string
	Source           neo4jentity.DataSource
	SourceOfTruth    neo4jentity.DataSource
	AppSource        string
	ParentID         string

	DataloaderKey string
}

type ServiceLineItemEntities []ServiceLineItemEntity
