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

type NoteEntities []NoteEntity

func (NoteEntity) IsTimelineEvent() {
}

func (NoteEntity) TimelineEventLabel() string {
	return NodeLabel_Note
}

func (note *NoteEntity) SetDataloaderKey(key string) {
	note.DataloaderKey = key
}

func (note NoteEntity) GetDataloaderKey() string {
	return note.DataloaderKey
}

func (note NoteEntity) Labels(tenant string) []string {
	return []string{"Note", "TimelineEvent", "Note_" + tenant, "TimelineEvent_" + tenant}
}
