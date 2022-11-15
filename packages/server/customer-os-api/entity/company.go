package entity

import (
	"fmt"
)

type CompanyEntity struct {
	Id   string
	Name string
}

func (company CompanyEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", company.Id, company.Name)
}

type CompanyEntities []CompanyEntity
