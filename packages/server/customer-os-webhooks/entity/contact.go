package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

// Deprecated, use neo4j module instead
type ContactEntity struct {
	Id              string
	Prefix          string
	Name            string
	FirstName       string
	LastName        string
	Description     string
	Timezone        string
	ProfilePhotoUrl string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Source          neo4jentity.DataSource
	SourceOfTruth   neo4jentity.DataSource
	AppSource       string
}

type ContactEntities []ContactEntity
