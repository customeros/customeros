package entity

import (
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string
	RawEmail      string
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EmailEntities []EmailEntity
