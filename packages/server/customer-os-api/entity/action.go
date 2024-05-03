package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

// Deprecated, use neo4j module instead
type ActionEntity struct {
	Id            string
	CreatedAt     time.Time
	Content       string
	Metadata      string
	Type          neo4jenum.ActionType
	Source        neo4jentity.DataSource
	AppSource     string
	DataloaderKey string
}

func (action ActionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", action.Id, action.Type)
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

type ActionEntities []ActionEntity

func (action ActionEntity) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelAction,
		neo4jutil.NodeLabelAction + "_" + tenant,
		neo4jutil.NodeLabelTimelineEvent,
		neo4jutil.NodeLabelTimelineEvent + "_" + tenant,
	}
}
