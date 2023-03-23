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

type AnalysisEntitys []AnalysisEntity

func (AnalysisEntity) IsTimelineEvent() {
}

func (AnalysisEntity) TimelineEventLabel() string {
	return NodeLabel_Analysis
}

func (AnalysisEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Analysis,
		NodeLabel_Analysis + "_" + tenant,
		NodeLabel_TimelineEvent,
		NodeLabel_TimelineEvent + "_" + tenant,
	}
}
