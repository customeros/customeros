package entity

import (
	"time"
)

type DomainEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Domain        string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type DomainEntities []DomainEntity
