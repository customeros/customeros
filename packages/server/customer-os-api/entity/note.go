package entity

import (
	"fmt"
	"time"
)

type NoteEntity struct {
	Id            string
	Html          string // deprecated
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (note NoteEntity) ToString() string {
	return fmt.Sprintf("id: %s\nhtml: %s", note.Id, note.Html)
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
