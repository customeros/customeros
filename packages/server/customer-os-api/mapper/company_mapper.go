package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapCompanyInputToEntity(input *model.CompanyInput) *entity.CompanyEntity {
	companyEntity := entity.CompanyEntity{}
	if input.ID != nil {
		companyEntity.Id = *input.ID
	}
	if input.Name != nil {
		companyEntity.Name = *input.Name
	}
	return &companyEntity
}

func MapEntityToCompany(entity *entity.CompanyEntity) *model.Company {
	return &model.Company{
		ID:   entity.Id,
		Name: entity.Name,
	}
}

//func MapEntitiesToCompanyPositiones(companyPositionEntities *entity.CompanyPositionEntities) []*model.Company {
//	var companyPositions []*model.Company
//	for _, companyPositionEntity := range *companyPositionEntities {
//		companyPositions = append(companyPositions, MapEntityToCompanyPosition(&companyPositionEntity))
//	}
//	return companyPositions
//}
