package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type NoteEntity struct {
	DataLoaderKey
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

func (NoteEntity) IsTimelineEvent() {
}

func (NoteEntity) TimelineEventLabel() string {
	return model.NodeLabelNote
}

func (e *NoteEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *NoteEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
