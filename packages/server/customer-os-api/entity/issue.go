package entity

import "time"

type IssueEntity struct {
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Subject     string
	Status      string
	Priority    string
	Description string

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

type IssueEntities []IssueEntity

func (*IssueEntity) IsTimelineEvent() {
}

func (*IssueEntity) TimelineEventLabel() string {
	return NodeLabel_Issue
}

func (issue *IssueEntity) SetDataloaderKey(key string) {
	issue.DataloaderKey = key
}

func (issue *IssueEntity) GetDataloaderKey() string {
	return issue.DataloaderKey
}

func (*IssueEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Issue,
		NodeLabel_Issue + "_" + tenant,
		NodeLabel_TimelineEvent,
		NodeLabel_TimelineEvent + "_" + tenant,
	}
}
