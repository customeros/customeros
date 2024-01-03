package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type UserEntity struct {
	Id               string
	FirstName        string
	LastName         string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Source           neo4jentity.DataSource
	SourceOfTruth    neo4jentity.DataSource
	AppSource        string
	Roles            []string
	Timezone         string
	ProfilePhotoUrl  string
	Internal         bool
	Bot              bool
	DefaultForPlayer bool
}

type UserEntities []UserEntity
