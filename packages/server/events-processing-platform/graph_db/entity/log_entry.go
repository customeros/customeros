package entity

import (
	"time"
)

type LogEntryEntity struct {
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
	return NodeLabel_LogEntry
}
