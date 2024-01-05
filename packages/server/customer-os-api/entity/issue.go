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
	return neo4jentity.NodeLabelIssue
}

func (issue *IssueEntity) SetDataloaderKey(key string) {
	issue.DataloaderKey = key
}

func (issue *IssueEntity) GetDataloaderKey() string {
	return issue.DataloaderKey
}

func (*IssueEntity) Labels(tenant string) []string {
	return []string{
		neo4jentity.NodeLabelIssue,
		neo4jentity.NodeLabelIssue + "_" + tenant,
		neo4jentity.NodeLabelTimelineEvent,
		neo4jentity.NodeLabelTimelineEvent + "_" + tenant,
	}
}
