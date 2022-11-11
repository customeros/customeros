package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

func MapCompanyPositionInputToEntity(input *model.CompanyInput) *entity.CompanyPositionEntity {
	if input == nil {
		return nil
	}
	companyPositionEntity := entity.CompanyPositionEntity{
		Company: input.CompanyName,
	}
	if input.JobTitle != nil {
		companyPositionEntity.JobTitle = *input.JobTitle
	}
	return &companyPositionEntity
}

func MapEntityToCompanyPosition(companyPosition *entity.CompanyPositionEntity) *model.Company {
	return &model.Company{
		CompanyName: companyPosition.Company,
		JobTitle:    utils.StringPtr(companyPosition.JobTitle),
	}
}

func MapEntitiesToCompanyPositiones(companyPositionEntities *entity.CompanyPositionEntities) []*model.Company {
	var companyPositions []*model.Company
	for _, companyPositionEntity := range *companyPositionEntities {
		companyPositions = append(companyPositions, MapEntityToCompanyPosition(&companyPositionEntity))
	}
	return companyPositions
}
