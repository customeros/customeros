package entity

import (
	"strings"
	"time"
)

type SocialProperty string

const (
	SocialPropertyId             SocialProperty = "id"
	SocialPropertyCreatedAt      SocialProperty = "createdAt"
	SocialPropertyUpdatedAt      SocialProperty = "updatedAt"
	SocialPropertySource         SocialProperty = "source"
	SocialPropertyAppSource      SocialProperty = "appSource"
	SocialPropertyUrl            SocialProperty = "url"
	SocialPropertyAlias          SocialProperty = "alias"
	SocialPropertyFollowersCount SocialProperty = "followersCount"
	SocialPropertyExternalId     SocialProperty = "externalId"
)

type SocialEntity struct {
	DataLoaderKey
	Id             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Source         DataSource
	AppSource      string
	Url            string
	Alias          string
	FollowersCount int64
	ExternalId     string
}

type SocialEntities []SocialEntity

func (s SocialEntity) IsLinkedin() bool {
	return strings.Contains(s.Url, "linkedin.com")
}
