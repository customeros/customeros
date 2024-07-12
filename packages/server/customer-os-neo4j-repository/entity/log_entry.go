package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
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
	return model.NodeLabelLogEntry
}

func (e *LogEntryEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *LogEntryEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
