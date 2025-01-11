package entity

import (
	"time"
)

type IndustryProperty string

const (
	IndustryPropertyCreatedAt IndustryProperty = "createdAt"
	IndustryPropertyCode      IndustryProperty = "code"
	IndustryPropertyName      IndustryProperty = "name"
)

type IndustryEntity struct {
	DataLoaderKey
	CreatedAt time.Time
	Code      string
	Name      string
}

type IndustryEntities []IndustryEntity
