package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type LogEntryEntity struct {
	EventStoreAggregate
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

func (LogEntryEntity) IsTimelineEvent() {
}

func (LogEntryEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelLogEntry
}
