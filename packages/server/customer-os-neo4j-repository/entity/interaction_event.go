package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type InteractionEventEntity struct {
	DataLoaderKey
	Id                           string
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	Content                      string
	ContentType                  string
	Channel                      string
	ChannelData                  string
	Identifier                   string
	CustomerOSInternalIdentifier string
	EventType                    string
	Hide                         bool
	Source                       DataSource
	SourceOfTruth                DataSource
	AppSource                    string
}

func (InteractionEventEntity) IsTimelineEvent() {
}

func (InteractionEventEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelInteractionEvent
}

func (e *InteractionEventEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *InteractionEventEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
