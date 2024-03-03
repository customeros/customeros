package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type IssueEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Subject       string
	Status        string
	Priority      string
	Description   string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type IssueEntities []IssueEntity

func (IssueEntity) IsTimelineEvent() {
}

func (IssueEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelIssue
}
