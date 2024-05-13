package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type InteractionEventEntity struct {
	Id                           string
	CreatedAt                    *time.Time
	Channel                      *string
	ChannelData                  *string
	ExternalId                   *string
	ExternalSystemId             *string
	EventIdentifier              string
	CustomerOSInternalIdentifier string
	Content                      string
	ContentType                  string
	Hide                         bool
	Source                       neo4jentity.DataSource
	SourceOfTruth                neo4jentity.DataSource
	EventType                    *string
	AppSource                    string

	DataloaderKey string
}

func (interactionEventEntity InteractionEventEntity) ToString() string {
	return fmt.Sprintf("id: %s", interactionEventEntity.Id)
}

type InteractionEventEntities []InteractionEventEntity

func (InteractionEventEntity) IsTimelineEvent() {
}

func (InteractionEventEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelInteractionEvent
}

func (InteractionEventEntity) IsAnalysisDescribe() {
}

func (InteractionEventEntity) AnalysisDescribeLabel() string {
	return neo4jutil.NodeLabelInteractionEvent
}

func (interactionEventEntity *InteractionEventEntity) SetDataloaderKey(key string) {
	interactionEventEntity.DataloaderKey = key
}

func (interactionEventEntity InteractionEventEntity) GetDataloaderKey() string {
	return interactionEventEntity.DataloaderKey
}

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelInteractionEvent,
		neo4jutil.NodeLabelInteractionEvent + "_" + tenant,
		neo4jutil.NodeLabelTimelineEvent,
		neo4jutil.NodeLabelTimelineEvent + "_" + tenant,
	}
}
