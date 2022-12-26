package entity

import (
	"fmt"
	"time"
)

type CompanyEntity struct {
	Id          string
	Name        string
	Description string
	Domain      string
	Website     string
	Industry    string
	IsPublic    bool
	CreatedAt   time.Time
}

func (company CompanyEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", company.Id, company.Name)
}

type CompanyEntities []CompanyEntity

func (company CompanyEntity) Labels() []string {
	return []string{"Company"}
}
