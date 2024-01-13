package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
}

func (ActionEntity) IsTimelineEvent() {
}

func (ActionEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelAction
}
