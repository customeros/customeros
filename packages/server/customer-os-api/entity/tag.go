package entity

import (
	"fmt"
	"time"
)

type TagEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	AppSource string
	TaggedAt  time.Time

	DataloaderKey string
}

func (tag TagEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", tag.Id, tag.Name)
}

type TagEntities []TagEntity

func (tag TagEntity) Labels(tenant string) []string {
	return []string{"Tag", "Tag_" + tenant}
}
