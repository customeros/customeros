package entity

type SourceFields struct {
	Source        DataSource `json:"source"`
	SourceOfTruth DataSource `json:"sourceOfTruth"`
	AppSource     string     `json:"appSource"`
}
