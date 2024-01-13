package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type IssueEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Subject       string
	Status        string
	Priority      string
	Description   string
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
}

func (IssueEntity) IsTimelineEvent() {
}

func (IssueEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelIssue
}
