package entity

import (
	"time"
)

type ContactEntity struct {
	DataLoaderKey
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	FirstName     string
	LastName      string
	Name          string
}

type ContactEntities []ContactEntity
