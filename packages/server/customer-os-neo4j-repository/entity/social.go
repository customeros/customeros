package entity

import "time"

type SocialProperty string

const (
	SocialPropertyId             SocialProperty = "id"
	SocialPropertyCreatedAt      SocialProperty = "createdAt"
	SocialPropertyUpdatedAt      SocialProperty = "updatedAt"
	SocialPropertySource         SocialProperty = "source"
	SocialPropertySourceOfTruth  SocialProperty = "sourceOfTruth"
	SocialPropertyAppSource      SocialProperty = "appSource"
	SocialPropertyUrl            SocialProperty = "url"
	SocialPropertyAlias          SocialProperty = "alias"
	SocialPropertyFollowersCount SocialProperty = "followersCount"
)

type SocialEntity struct {
	DataLoaderKey
	Id             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Source         DataSource
	SourceOfTruth  DataSource
	AppSource      string
	Url            string
	Alias          string
	FollowersCount int64
}

type SocialEntities []SocialEntity
