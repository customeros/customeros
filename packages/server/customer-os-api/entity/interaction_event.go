package entity

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
	return model.NodeLabelInteractionEvent
}

func (InteractionEventEntity) IsAnalysisDescribe() {
}

func (InteractionEventEntity) AnalysisDescribeLabel() string {
	return model.NodeLabelInteractionEvent
}

func (interactionEventEntity *InteractionEventEntity) SetDataloaderKey(key string) {
	interactionEventEntity.DataloaderKey = key
}

func (interactionEventEntity InteractionEventEntity) GetDataloaderKey() string {
	return interactionEventEntity.DataloaderKey
}

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelInteractionEvent,
		model.NodeLabelInteractionEvent + "_" + tenant,
		model.NodeLabelTimelineEvent,
		model.NodeLabelTimelineEvent + "_" + tenant,
	}
}
