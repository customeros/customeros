package entity

import (
	"time"
)

type UserEntity struct {
	Id               string
	FirstName        string
	LastName         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	Roles            []string
	ProfilePhotoUrl  string
	Timezone         string
	Internal         bool
	Bot              bool
	DefaultForPlayer bool
	Tenant           string
}

type UserEntities []UserEntity
