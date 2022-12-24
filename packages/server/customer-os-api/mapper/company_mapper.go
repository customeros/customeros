package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapCompanyInputToEntity(input *model.CompanyInput) *entity.CompanyEntity {
	return &entity.CompanyEntity{
		Id:   utils.IfNotNilString(input.ID),
		Name: utils.IfNotNilString(input.Name),
	}
}

func MapEntityToCompany(entity *entity.CompanyEntity) *model.Company {
	return &model.Company{
		ID:          entity.Id,
		Name:        entity.Name,
		Description: utils.StringPtr(entity.Description),
		Domain:      utils.StringPtr(entity.Domain),
		Website:     utils.StringPtr(entity.Website),
		Industry:    utils.StringPtr(entity.Industry),
		IsPublic:    utils.BoolPtr(entity.IsPublic),
		CreatedAt:   entity.CreatedAt,
	}
}

func MapEntitiesToCompanies(companyEntities *entity.CompanyEntities) []*model.Company {
	var companies []*model.Company
	for _, companyEntity := range *companyEntities {
		companies = append(companies, MapEntityToCompany(&companyEntity))
	}
	return companies
}
