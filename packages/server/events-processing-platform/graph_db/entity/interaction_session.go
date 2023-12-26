package entity

import (
	"time"
)

type InteractionSessionEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Identifier    string
	Name          string
	Status        string
	Type          string
	Channel       string
	ChannelData   string
	AppSource     string
	Source        DataSource
	SourceOfTruth DataSource
}

func (InteractionSessionEntity) IsTimelineEvent() {
}

func (InteractionSessionEntity) TimelineEventLabel() string {
	return NodeLabel_InteractionSession
}
