package entity

import "time"

type InteractionEventEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Content       string
	ContentType   string
	Channel       string
	ChannelData   string
	Identifier    string
	EventType     string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type InteractionEventEntities []InteractionEventEntity
