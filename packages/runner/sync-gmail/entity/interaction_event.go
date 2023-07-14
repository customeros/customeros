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
	Source          string
	SourceOfTruth   string
	EventType       *string
	AppSource       string
}

func (interactionEventEntity InteractionEventEntity) ToString() string {
	return fmt.Sprintf("id: %s", interactionEventEntity.Id)
}

func (InteractionEventEntity) Labels(tenant string) []string {
	return []string{"InteractionEvent", "TimelineEvent", "InteractionEvent_" + tenant, "TimelineEvent_" + tenant}
}
