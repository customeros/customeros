package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ActionEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	Content       string
	Metadata      string
	Type          neo4jenum.ActionType
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type ActionEntities []ActionEntity

func (e *ActionEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *ActionEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}

func (ActionEntity) IsTimelineEvent() {
}

func (ActionEntity) TimelineEventLabel() string {
	return model.NodeLabelAction
}

func (e *ActionEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelAction,
		model.NodeLabelAction + "_" + tenant,
		model.NodeLabelTimelineEvent,
		model.NodeLabelTimelineEvent + "_" + tenant,
	}
}
