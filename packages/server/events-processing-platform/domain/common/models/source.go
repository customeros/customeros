package models

import "fmt"

type Source struct {
	Source        string `json:"source"`
	SourceOfTruth string `json:"sourceOfTruth"`
	AppSource     string `json:"appSource"`
}

func (s *Source) String() string {
	return fmt.Sprintf("Source{Source: %s, SourceOfTruth: %s, AppSource: %s}", s.Source, s.SourceOfTruth, s.AppSource)
}
