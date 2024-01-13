package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
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
	return neo4jutil.NodeLabelIssue
}

func (issue *IssueEntity) SetDataloaderKey(key string) {
	issue.DataloaderKey = key
}

func (issue *IssueEntity) GetDataloaderKey() string {
	return issue.DataloaderKey
}

func (*IssueEntity) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelIssue,
		neo4jutil.NodeLabelIssue + "_" + tenant,
		neo4jutil.NodeLabelTimelineEvent,
		neo4jutil.NodeLabelTimelineEvent + "_" + tenant,
	}
}
