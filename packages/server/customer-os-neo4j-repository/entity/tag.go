package entity

import (
	"time"
)

type TagProperty string

const (
	TagPropertyId   TagProperty = "id"
	TagPropertyName TagProperty = "name"
)

type TagEntity struct {
	DataLoaderKey
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	TaggedAt  time.Time
}

type TagEntities []TagEntity
