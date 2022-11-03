package entity

import (
	"fmt"
)

type ContactCompanyEntity struct {
	Company  string
	JobTitle string
}

func (contactCompany ContactCompanyEntity) ToString() string {
	return fmt.Sprintf("company: %s\njob title: %s", contactCompany.Company, contactCompany.JobTitle)
}

type ContactCompanyEntities []ContactCompanyEntity
