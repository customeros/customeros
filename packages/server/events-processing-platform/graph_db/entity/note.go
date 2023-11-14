package entity

import (
	"time"
)

type NoteEntity struct {
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (NoteEntity) IsTimelineEvent() {
}

func (NoteEntity) TimelineEventLabel() string {
	return NodeLabel_Note
}
