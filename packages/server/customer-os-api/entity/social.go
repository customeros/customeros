package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
		neo4jentity.NodeLabel_Social,
		neo4jentity.NodeLabel_Social + "_" + tenant,
	}
}
