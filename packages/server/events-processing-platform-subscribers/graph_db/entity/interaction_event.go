package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

// Deprecated
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
	Hide          bool
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
}

func (InteractionEventEntity) IsTimelineEvent() {
}

func (InteractionEventEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelInteractionEvent
}
