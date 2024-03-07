package entity

import (
	"time"
)

type CommentEntity struct {
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}
