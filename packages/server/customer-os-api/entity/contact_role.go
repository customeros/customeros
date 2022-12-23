package entity

import (
	"fmt"
)

type ContactRoleEntity struct {
	Id       string
	JobTitle string
	Primary  bool
}

func (contactRole ContactRoleEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", contactRole.Id, contactRole.JobTitle)
}

type ContactRoleEntities []ContactRoleEntity
