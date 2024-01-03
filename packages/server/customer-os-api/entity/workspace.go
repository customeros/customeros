package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type WorkspaceEntity struct {
	Id            string
	Name          string
	Provider      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

func (workspace WorkspaceEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", workspace.Id, workspace.Name)
}

type WorkspaceEntities []WorkspaceEntity

func (workspace WorkspaceEntity) Labels(tenant string) []string {
	return []string{"Domain"}
}
