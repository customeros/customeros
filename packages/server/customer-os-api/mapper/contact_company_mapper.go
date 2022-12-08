package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapCompanyPositionInputToEntity(input *model.CompanyPositionInput) *entity.CompanyPositionEntity {
	if input == nil {
		return nil
	}
	companyPositionEntity := entity.CompanyPositionEntity{
		Company:  *MapCompanyInputToEntity(input.Company),
		JobTitle: utils.IfNotNilString(input.JobTitle),
	}
	return &companyPositionEntity
}

func MapEntityToCompanyPosition(entity *entity.CompanyPositionEntity) *model.CompanyPosition {
	return &model.CompanyPosition{
		ID:       entity.Id,
		JobTitle: utils.StringPtr(entity.JobTitle),
		Company:  MapEntityToCompany(&entity.Company),
	}
}

func MapEntitiesToCompanyPositiones(companyPositionEntities *entity.CompanyPositionEntities) []*model.CompanyPosition {
	var companyPositions []*model.CompanyPosition
	for _, companyPositionEntity := range *companyPositionEntities {
		companyPositions = append(companyPositions, MapEntityToCompanyPosition(&companyPositionEntity))
	}
	return companyPositions
}
