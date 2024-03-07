package entity

import (
	"time"
)

type SocialEntity struct {
	DataLoaderKey
	Id            string
	PlatformName  string
	Url           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type SocialEntities []SocialEntity
