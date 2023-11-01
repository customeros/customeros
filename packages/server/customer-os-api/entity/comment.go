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

	DataloaderKey string
}

type CommentEntities []CommentEntity

func (comment *CommentEntity) SetDataloaderKey(key string) {
	comment.DataloaderKey = key
}

func (comment *CommentEntity) GetDataloaderKey() string {
	return comment.DataloaderKey
}

func (CommentEntity) Labels(tenant string) []string {
	return []string{NodeLabel_Comment, NodeLabel_Comment + "_" + tenant}
}
