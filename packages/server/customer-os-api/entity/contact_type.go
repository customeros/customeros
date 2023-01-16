package entity

import "fmt"

type ContactTypeEntity struct {
	Id   string
	Name string
}

func (contactType ContactTypeEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactType.Id, contactType.Name)
}

type ContactTypeEntities []ContactTypeEntity

func (contactType ContactTypeEntity) Labels() []string {
	return []string{"ContactType"}
}
