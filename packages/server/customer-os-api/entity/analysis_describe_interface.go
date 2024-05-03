package entity

type AnalysisDescribeDetails struct {
	Type string
}

// Deprecated, use neo4j module instead
type AnalysisDescribe interface {
	IsAnalysisDescribe()
	AnalysisDescribeLabel() string
	GetDataloaderKey() string
}

type AnalysisDescribes []AnalysisDescribe
