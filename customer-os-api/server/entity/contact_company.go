package entity

import (
	"fmt"
)

type CompanyPositionEntity struct {
	Company  string
	JobTitle string
}

func (companyPosition CompanyPositionEntity) ToString() string {
	return fmt.Sprintf("company: %s\njob title: %s", companyPosition.Company, companyPosition.JobTitle)
}

type CompanyPositionEntities []CompanyPositionEntity
