package entity

type SearchResultEntityType string

const (
	SearchResultEntityTypeContact      SearchResultEntityType = "CONTACT"
	SearchResultEntityTypeOrganization SearchResultEntityType = "ORGANIZATION"
	SearchResultEntityTypeEmail        SearchResultEntityType = "EMAIL"
	SearchResultEntityTypeState        SearchResultEntityType = "STATE"
)

type SearchResultEntity struct {
	Score      float64
	Labels     []string
	Node       any
	EntityType SearchResultEntityType
}

type SearchResultEntities []SearchResultEntity
