package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type LogEntryEntity struct {
	EventStoreAggregate
	DataLoaderKey
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type LogEntryEntities []LogEntryEntity

func (LogEntryEntity) IsTimelineEvent() {
}

func (LogEntryEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelLogEntry
}

func (e *LogEntryEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *LogEntryEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
