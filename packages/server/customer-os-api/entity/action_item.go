package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type ActionItemEntity struct {
	Id        string
	CreatedAt *time.Time

	Content string

	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

func (entity ActionItemEntity) ToString() string {
	return fmt.Sprintf("id: %s", entity.Id)
}

type ActionItemEntities []ActionItemEntity

func (entity ActionItemEntity) Labels(tenant string) []string {
	return []string{
		neo4jentity.NodeLabel_ActionItem,
		neo4jentity.NodeLabel_ActionItem + "_" + tenant,
	}
}
