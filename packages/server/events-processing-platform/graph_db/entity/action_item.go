package entity

import (
	"fmt"
	"time"
)

type ActionItemEntity struct {
	Id        string
	CreatedAt *time.Time

	Content string

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (entity ActionItemEntity) ToString() string {
	return fmt.Sprintf("id: %s", entity.Id)
}

type ActionItemEntities []ActionItemEntity

func (entity ActionItemEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_ActionItem,
		NodeLabel_ActionItem + "_" + tenant,
	}
}
