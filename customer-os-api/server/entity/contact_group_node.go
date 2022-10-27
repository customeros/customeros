package entity

import (
	"fmt"
)

type ContactGroupNode struct {
	Id   string
	Name string
}

func (contactGroup ContactGroupNode) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactGroup.Id, contactGroup.Name)
}

type ContactGroupNodes []ContactGroupNode
