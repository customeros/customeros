package entity

import (
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

func (AnalysisEntity) IsTimelineEvent() {
}

func (AnalysisEntity) TimelineEventLabel() string {
	return NodeLabel_Analysis
}
