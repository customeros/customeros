package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

// Deprecated: use neo4jentity.TagEntity instead
type TagEntity struct {
	Id            string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
	TaggedAt      time.Time

	DataloaderKey string
}

func (tag TagEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", tag.Id, tag.Name)
}

// Deprecated: use neo4jentity.TagEntities instead
type TagEntities []TagEntity

func (TagEntity) Labels(tenant string) []string {
	return []string{
		neo4jutil.NodeLabelTag,
		neo4jutil.NodeLabelTag + "_" + tenant,
	}
}
