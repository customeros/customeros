package entity

import (
	"fmt"
	"time"
)

type WorkspaceEntity struct {
	Id            string
	Name          string
	Provider      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
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
