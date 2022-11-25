package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapCompanyInputToEntity(input *model.CompanyInput) *entity.CompanyEntity {
	companyEntity := new(entity.CompanyEntity)
	if input.ID != nil {
		companyEntity.Id = *input.ID
	}
	if input.Name != nil {
		companyEntity.Name = *input.Name
	}
	return companyEntity
}

func MapEntityToCompany(entity *entity.CompanyEntity) *model.Company {
	return &model.Company{
		ID:   entity.Id,
		Name: entity.Name,
	}
}

func MapEntitiesToCompanies(companyEntities *entity.CompanyEntities) []*model.Company {
	var companies []*model.Company
	for _, companyEntity := range *companyEntities {
		companies = append(companies, MapEntityToCompany(&companyEntity))
	}
	return companies
}
