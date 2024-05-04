package entity

import "time"

type SocialEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Url           string
}

type SocialEntities []SocialEntity
