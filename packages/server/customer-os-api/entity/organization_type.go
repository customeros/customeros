package entity

import (
	"fmt"
	"time"
)

type OrganizationTypeEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (organizationType OrganizationTypeEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organizationType.Id, organizationType.Name)
}

type OrganizationTypeEntities []OrganizationTypeEntity

func (organizationType OrganizationTypeEntity) Labels() []string {
	return []string{"OrganizationType"}
}
