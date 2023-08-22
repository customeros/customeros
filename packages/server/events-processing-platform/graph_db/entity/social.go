package entity

import (
	"fmt"
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

func (social SocialEntity) ToString() string {
	return fmt.Sprintf("id: %s paltform: %s url: %s", social.Id, social.PlatformName, social.Url)
}

type SocialEntities []SocialEntity

func (social SocialEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Social,
		NodeLabel_Social + "_" + tenant,
	}
}
