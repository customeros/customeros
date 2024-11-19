package entity

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type TagProperty string

const (
	TagPropertyId         TagProperty = "id"
	TagPropertyName       TagProperty = "name"
	TagPropertyEntityType TagProperty = "entityType"
	TagPropertySource     TagProperty = "source"
	TagPropertyAppSource  TagProperty = "appSource"
	TagPropertyCreatedAt  TagProperty = "createdAt"
	TagPropertyUpdatedAt  TagProperty = "updatedAt"
)

type TagEntity struct {
	DataLoaderKey
	Id         string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Source     DataSource
	AppSource  string
	TaggedAt   time.Time
	EntityType commonmodel.EntityType
}

type TagEntities []TagEntity
