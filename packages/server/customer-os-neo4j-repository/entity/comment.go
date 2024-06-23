package entity

import (
	"time"
)

type CommentEntity struct {
	DataloaderKey string
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type CommentEntities []CommentEntity
