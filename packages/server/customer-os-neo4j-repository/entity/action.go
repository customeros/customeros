package entity

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type ActionEntity struct {
	Id            string
	CreatedAt     time.Time
	Content       string
	Metadata      string
	Type          neo4jenum.ActionType
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (action ActionEntity) GetDataloaderKey() string {
	return action.DataloaderKey
}

func (action *ActionEntity) SetDataloaderKey(key string) {
	action.DataloaderKey = key
}

func (ActionEntity) IsTimelineEvent() {
}

func (ActionEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelAction
}
