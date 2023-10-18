package entity

import (
	"time"
)

type SocialEntity struct {
	Id           string
	PlatformName string
	Url          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SourceFields SourceFields

	DataloaderKey string
}

type SocialEntities []SocialEntity
