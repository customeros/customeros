package model

type Source struct {
	Source        string `json:"source"`
	SourceOfTruth string `json:"sourceOfTruth"`
	AppSource     string `json:"appSource"`
}
