package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type IssueEntity struct {
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Subject     string
	Status      string
	Priority    string
	Description string

	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

type IssueEntities []IssueEntity

func (*IssueEntity) IsTimelineEvent() {
}

func (*IssueEntity) TimelineEventLabel() string {
	return neo4jentity.NodeLabel_Issue
}

func (issue *IssueEntity) SetDataloaderKey(key string) {
	issue.DataloaderKey = key
}

func (issue *IssueEntity) GetDataloaderKey() string {
	return issue.DataloaderKey
}

func (*IssueEntity) Labels(tenant string) []string {
	return []string{
		neo4jentity.NodeLabel_Issue,
		neo4jentity.NodeLabel_Issue + "_" + tenant,
		neo4jentity.NodeLabel_TimelineEvent,
		neo4jentity.NodeLabel_TimelineEvent + "_" + tenant,
	}
}
