package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type NoteEntity struct {
	DataLoaderKey string
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
	return neo4jutil.NodeLabelNote
}
