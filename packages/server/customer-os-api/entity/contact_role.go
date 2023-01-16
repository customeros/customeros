package entity

import (
	"fmt"
)

type ContactRoleEntity struct {
	Id            string
	JobTitle      string
	Primary       bool
	Source        DataSource
	SourceOfTruth DataSource
}

func (contactRole ContactRoleEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", contactRole.Id, contactRole.JobTitle)
}

type ContactRoleEntities []ContactRoleEntity

func (contactRole ContactRoleEntity) Labels() []string {
	return []string{"Role"}
}
