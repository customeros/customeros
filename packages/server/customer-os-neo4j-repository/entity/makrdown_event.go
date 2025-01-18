package neo4j_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"time"
)

type MarkdownEventEntity struct {
	DataLoaderKey
	Id        string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	AppSource string
}

type MarkdownEventEntities []MarkdownEventEntity

func (MarkdownEventEntity) IsTimelineEvent() {
}

func (MarkdownEventEntity) TimelineEventLabel() string {
	return model.NodeLabelMarkdownEvent
}

func (e *MarkdownEventEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *MarkdownEventEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
