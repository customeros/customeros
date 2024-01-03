package entity

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
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
}

func (LogEntryEntity) IsTimelineEvent() {
}

func (LogEntryEntity) TimelineEventLabel() string {
	return NodeLabel_LogEntry
}
