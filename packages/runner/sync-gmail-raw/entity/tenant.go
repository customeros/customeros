package entity

import (
	"fmt"
)

type TenantEntity struct {
	Id   string
	Name string
}

func (domain TenantEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", domain.Id, domain.Name)
}
