package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"time"
)

type InteractionEventEntity struct {
	Id               string
	CreatedAt        *time.Time
	Channel          *string
	ChannelData      *string
	ExternalId       *string
	ExternalSystemId *string
	EventIdentifier  string
	Content          string
	ContentType      string
	Hide             bool
	Source           neo4jentity.DataSource
	SourceOfTruth    neo4jentity.DataSource
	EventType        *string
	AppSource        string

	DataloaderKey string
}

func (interactionEventEntity InteractionEventEntity) ToString() string {
	return fmt.Sprintf("id: %s", interactionEventEntity.Id)
}

type InteractionEventEntities []InteractionEventEntity

func (InteractionEventEntity) IsTimelineEvent() {
}

func (InteractionEventEntity) TimelineEventLabel() string {
	return NodeLabel_InteractionEvent
}

func (InteractionEventEntity) IsAnalysisDescribe() {
}

func (InteractionEventEntity) AnalysisDescribeLabel() string {
	return NodeLabel_InteractionEvent
}

func (interactionEventEntity *InteractionEventEntity) SetDataloaderKey(key string) {
	interactionEventEntity.DataloaderKey = key
}

func (interactionEventEntity InteractionEventEntity) GetDataloaderKey() string {
	return interactionEventEntity.DataloaderKey
}

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_InteractionEvent,
		NodeLabel_InteractionEvent + "_" + tenant,
		NodeLabel_TimelineEvent,
		NodeLabel_TimelineEvent + "_" + tenant,
	}
}
