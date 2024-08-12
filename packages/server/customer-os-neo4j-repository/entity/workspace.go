package entity

import (
	"time"
)

type WorkspaceEntity struct {
	Id            string
	Name          string
	Provider      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

type WorkspaceEntities []WorkspaceEntity
