package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type DomainEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Domain        string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	EnrichDetails DomainEnrichDetails
}

type DomainEnrichDetails struct {
	EnrichRequestedAt *time.Time
	EnrichError       string
	EnrichedAt        *time.Time
	EnrichSource      enum.DomainEnrichSource
	EnrichData        string
}

type DomainEntities []DomainEntity
