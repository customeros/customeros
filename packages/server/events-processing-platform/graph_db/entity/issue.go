package entity

import "time"

type IssueEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Subject       string
	Status        string
	Priority      string
	Description   string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type IssueEntities []IssueEntity
