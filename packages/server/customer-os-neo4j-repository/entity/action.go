package entity

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
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
	return neo4jutil.NodeLabelAction
}

func (e *ActionEntity) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelAction,
		neo4jutil.NodeLabelAction + "_" + tenant,
		neo4jutil.NodeLabelTimelineEvent,
		neo4jutil.NodeLabelTimelineEvent + "_" + tenant,
	}
}
