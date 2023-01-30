package entity

import (
	"fmt"
	"time"
)

type ContactRoleEntity struct {
	Id                  string
	JobTitle            string
	Primary             bool
	ResponsibilityLevel int64
	Source              DataSource
	SourceOfTruth       DataSource
	AppSource           string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (contactRole ContactRoleEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", contactRole.Id, contactRole.JobTitle)
}

type ContactRoleEntities []ContactRoleEntity

func (contactRole ContactRoleEntity) Labels() []string {
	return []string{"Role"}
}
