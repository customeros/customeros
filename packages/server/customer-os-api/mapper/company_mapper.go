package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapCompanyInputToEntity(input *model.CompanyInput) *entity.CompanyEntity {
	return &entity.CompanyEntity{
		Name:        input.Name,
		Description: utils.IfNotNilString(input.Description),
		Domain:      utils.IfNotNilString(input.Domain),
		Website:     utils.IfNotNilString(input.Website),
		Industry:    utils.IfNotNilString(input.Industry),
		IsPublic:    utils.IfNotNilBool(input.IsPublic),
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
		Readonly:    utils.BoolPtr(entity.Readonly),
	}
}

func MapEntitiesToCompanies(companyEntities *entity.CompanyEntities) []*model.Company {
	var companies []*model.Company
	for _, companyEntity := range *companyEntities {
		companies = append(companies, MapEntityToCompany(&companyEntity))
	}
	return companies
}
