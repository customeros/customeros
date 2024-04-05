package entity

import "time"

type SocialEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	Url string

	DataLoaderKey
}

type SocialEntities []SocialEntity
