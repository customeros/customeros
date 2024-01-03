package entity

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"time"
)

type NoteEntity struct {
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

func (NoteEntity) IsTimelineEvent() {
}

func (NoteEntity) TimelineEventLabel() string {
	return NodeLabel_Note
}
