package entity

import (
	"fmt"
	"time"
)

type InteractionEventEntity struct {
	Id              string
	CreatedAt       *time.Time
	Channel         string
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

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_InteractionEvent,
		NodeLabel_InteractionEvent + "_" + tenant,
		NodeLabel_Action,
		NodeLabel_Action + "_" + tenant,
	}
}
