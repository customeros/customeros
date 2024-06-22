package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type AnalysisEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	Content       string
	ContentType   string
	AnalysisType  string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

func (AnalysisEntity) IsTimelineEvent() {
}

func (AnalysisEntity) TimelineEventLabel() string {
	return neo4jutil.NodeLabelAnalysis
}

func (e *AnalysisEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *AnalysisEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}
