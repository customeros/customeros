package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type LogEntryEntity struct {
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

type LogEntryEntities []LogEntryEntity

func (*LogEntryEntity) IsTimelineEvent() {
}

func (*LogEntryEntity) TimelineEventLabel() string {
	return NodeLabel_LogEntry
}

func (logEntry *LogEntryEntity) SetDataloaderKey(key string) {
	logEntry.DataloaderKey = key
}

func (logEntry *LogEntryEntity) GetDataloaderKey() string {
	return logEntry.DataloaderKey
}

func (LogEntryEntity) Labels(tenant string) []string {
	return []string{"LogEntry", "TimelineEvent", "LogEntry_" + tenant, "TimelineEvent_" + tenant}
}
