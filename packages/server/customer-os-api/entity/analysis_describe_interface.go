package entity

type AnalysisDescribeDetails struct {
	Type string
}

type AnalysisDescribe interface {
	IsAnalysisDescribe()
	AnalysisDescribeLabel() string
	GetDataloaderKey() string
}

type AnalysisDescribes []AnalysisDescribe
