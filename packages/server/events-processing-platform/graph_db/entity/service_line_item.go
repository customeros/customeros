package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

// Deprecated
type ServiceLineItemEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	EndedAt       *time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
	Name          string
	Billed        string
	Price         float64
	Quantity      int64
	Comments      string
	ParentId      string
	IsCanceled    bool
}
