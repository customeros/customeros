package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
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
	return neo4jutil.NodeLabelAnalysis
}
