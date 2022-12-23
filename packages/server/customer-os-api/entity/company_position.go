package entity

import (
	"fmt"
)

// TODO alexb delete
type CompanyPositionEntity struct {
	Id       string
	Company  CompanyEntity
	JobTitle string
}

func (companyPosition CompanyPositionEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", companyPosition.Id, companyPosition.JobTitle)
}

type CompanyPositionEntities []CompanyPositionEntity
