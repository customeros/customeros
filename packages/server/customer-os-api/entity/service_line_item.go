package entity

import (
	"time"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

// Deprecated
type ServiceLineItemEntity struct {
	ID               string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	StartedAt        time.Time
	EndedAt          *time.Time
	IsCanceled       bool
	Billed           neo4jenum.BilledType
	Price            float64
	Quantity         int64
	PreviousBilled   neo4jenum.BilledType
	PreviousPrice    float64
	PreviousQuantity int64
	Comments         string
	Source           neo4jentity.DataSource
	SourceOfTruth    neo4jentity.DataSource
	AppSource        string
	ParentID         string

	DataloaderKey string
}

// Deprecated
type ServiceLineItemEntities []ServiceLineItemEntity
