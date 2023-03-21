package entity

import (
	"fmt"
	"time"
)

type InteractionSessionEntity struct {
	Id                string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EndedAt           *time.Time
	SessionIdentifier *string
	Name              string
	Status            string
	Type              *string
	Channel           *string
	ChannelData       *string
	AppSource         string
	Source            DataSource
	SourceOfTruth     DataSource

	DataloaderKey string
}

func (interactionSession InteractionSessionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", interactionSession.Id, interactionSession.Name)
}

type InteractionSessionEntities []InteractionSessionEntity

func (InteractionSessionEntity) IsTimelineEvent() {
}

func (InteractionSessionEntity) TimelineEventLabel() string {
	return NodeLabel_InteractionSession
}

func (interactionSession InteractionSessionEntity) Labels(tenant string) []string {
	return []string{"InteractionSession", "TimelineEvent", "InteractionSession_" + tenant, "TimelineEvent_" + tenant}
}
