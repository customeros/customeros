package entity

import (
	"fmt"
	"time"
)

type InteractionSessionEntity struct {
	Id            string
	StartedAt     time.Time
	EndedAt       *time.Time
	Name          string
	Status        string
	Type          string
	Channel       string
	AppSource     string
	Source        DataSource
	SourceOfTruth DataSource
}

func (interactionSession InteractionSessionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", interactionSession.Id, interactionSession.Name)
}

type InteractionSessionEntities []InteractionSessionEntity

func (InteractionSessionEntity) Action() {
}

func (InteractionSessionEntity) ActionName() string {
	return NodeLabel_InteractionSession
}

func (interactionSession InteractionSessionEntity) Labels(tenant string) []string {
	return []string{"InteractionSession", "Action", "InteractionSession_" + tenant, "Action_" + tenant}
}
