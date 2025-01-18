package neo4j_entity

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
	Code      string // NAICS code, see https://www.naics.com/search/
	Name      string // NAICS name
}

type IndustryEntities []IndustryEntity
