package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

// Deprecated, use neo4j module instead
type InteractionEventEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Content       string
	ContentType   string
	Channel       string
	ChannelData   string
	Identifier    string
	EventType     string
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
}

type InteractionEventEntities []InteractionEventEntity
