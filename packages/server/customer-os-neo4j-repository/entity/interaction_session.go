package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type InteractionSessionEntity struct {
	DataLoaderKey
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
	return neo4jutil.NodeLabelInteractionSession
}

func (e *InteractionSessionEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *InteractionSessionEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
