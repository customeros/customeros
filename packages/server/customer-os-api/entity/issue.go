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
}

func (IssueEntity) IsTimelineEvent() {
}

func (IssueEntity) TimelineEventLabel() string {
	return NodeLabel_Issue
}

func (IssueEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Issue,
		NodeLabel_Issue + "_" + tenant,
		NodeLabel_TimelineEvent,
		NodeLabel_TimelineEvent + "_" + tenant,
	}
}
