package entity

import (
	"time"
)

type UserEntity struct {
	Id               string
	FirstName        string
	LastName         string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	Roles            []string
	Timezone         string
	ProfilePhotoUrl  string
	Internal         bool
	Bot              bool
	DefaultForPlayer bool
}

type UserEntities []UserEntity
