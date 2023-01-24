package entity

import (
	"fmt"
	"time"
)

type ContactTypeEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

func (contactType ContactTypeEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactType.Id, contactType.Name)
}

type ContactTypeEntities []ContactTypeEntity

func (contactType ContactTypeEntity) Labels() []string {
	return []string{"ContactType"}
}
