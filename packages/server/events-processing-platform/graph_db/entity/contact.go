package entity

import (
	"time"
)

type ContactEntity struct {
	Id              string
	Prefix          string
	Name            string
	FirstName       string
	LastName        string
	Description     string
	Timezone        string
	ProfilePhotoUrl string
	CreatedAt       *time.Time
	UpdatedAt       time.Time
	Source          DataSource
	SourceOfTruth   DataSource
	AppSource       string
}

type ContactEntities []ContactEntity
