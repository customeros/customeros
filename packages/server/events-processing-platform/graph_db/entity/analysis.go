package entity

import (
	"fmt"
	"time"
)

type AnalysisEntity struct {
	Id        string
	CreatedAt *time.Time

	Content       string
	ContentType   string
	AnalysisType  string
	Source        DataSource
	SourceOfTruth DataSource
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
