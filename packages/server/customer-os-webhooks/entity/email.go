package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string
	RawEmail      string
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EmailEntities []EmailEntity
