package entity

import (
	"fmt"
	"time"
)

type InteractionEventEntity struct {
	Id              string
	CreatedAt       time.Time
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

func (InteractionEventEntity) TimelineEvent() {
}

func (InteractionEventEntity) TimelineEventName() string {
	return NodeLabel_InteractionEvent
}

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{"InteractionEvent", "InteractionEvent_" + tenant}
}
