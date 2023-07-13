package entity

import (
	"fmt"
	"time"
)

type InteractionSessionEntity struct {
	Id                string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	SessionIdentifier *string
	Name              string
	Status            string
	Type              *string
	Channel           *string
	ChannelData       *string
	AppSource         string
	Source            string
	SourceOfTruth     string

	DataloaderKey string
}

func (interactionSession InteractionSessionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", interactionSession.Id, interactionSession.Name)
}

func (interactionSession InteractionSessionEntity) Labels(tenant string) []string {
	return []string{"InteractionSession", "TimelineEvent", "InteractionSession_" + tenant, "TimelineEvent_" + tenant}
}
