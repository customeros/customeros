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

func (analysis *AnalysisEntity) SetDataloaderKey(key string) {
	analysis.DataloaderKey = key
}

func (analysis AnalysisEntity) GetDataloaderKey() string {
	return analysis.DataloaderKey
}

func (AnalysisEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Analysis,
		NodeLabel_Analysis + "_" + tenant,
	}
}
