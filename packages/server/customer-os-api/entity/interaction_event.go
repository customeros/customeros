package entity

import (
	"fmt"
	"time"
)

type InteractionEventEntity struct {
	Id              string
	CreatedAt       *time.Time
	Channel         *string
	ChannelData     *string
	EventIdentifier string
	Content         string
	ContentType     string
	Source          DataSource
	SourceOfTruth   DataSource
	AppSource       string

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
