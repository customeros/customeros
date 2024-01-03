package entity

import (
	"fmt"
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

func (analysisEntity AnalysisEntity) ToString() string {
	return fmt.Sprintf("id: %s", analysisEntity.Id)
}

type AnalysisEntities []AnalysisEntity

func (AnalysisEntity) IsTimelineEvent() {
}

func (analysis *AnalysisEntity) SetDataloaderKey(key string) {
	analysis.DataloaderKey = key
}

func (analysis AnalysisEntity) GetDataloaderKey() string {
	return analysis.DataloaderKey
}

func (AnalysisEntity) TimelineEventLabel() string {
	return NodeLabel_Analysis
}

func (AnalysisEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Analysis,
		NodeLabel_Analysis + "_" + tenant,
	}
}
