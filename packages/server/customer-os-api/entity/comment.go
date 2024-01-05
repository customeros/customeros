package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type CommentEntity struct {
	Id            string
	Content       string
	ContentType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
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
	return []string{neo4jentity.NodeLabelComment, neo4jentity.NodeLabelComment + "_" + tenant}
}
