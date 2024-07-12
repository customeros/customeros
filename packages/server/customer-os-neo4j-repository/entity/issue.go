package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
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
	return model.NodeLabelIssue
}

func (e *IssueEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *IssueEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
