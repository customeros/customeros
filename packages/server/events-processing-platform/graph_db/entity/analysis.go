package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type AnalysisEntity struct {
	Id        string
	CreatedAt *time.Time

	Content       string
	ContentType   string
	AnalysisType  string
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

func (AnalysisEntity) IsTimelineEvent() {
}

func (AnalysisEntity) TimelineEventLabel() string {
	return NodeLabel_Analysis
}
