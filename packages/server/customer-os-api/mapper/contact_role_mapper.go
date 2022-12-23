package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

//func MapContactRoleInputToEntity(input *model.ContactRoleInput) *entity.CompanyPositionEntity {
//	if input == nil {
//		return nil
//	}
//	companyPositionEntity := entity.CompanyPositionEntity{
//		Company:  *MapCompanyInputToEntity(input.Company),
//		JobTitle: utils.IfNotNilString(input.JobTitle),
//	}
//	return &companyPositionEntity
//}

func MapEntityToContactRole(entity *entity.ContactRoleEntity) *model.ContactRole {
	contactRole := model.ContactRole{
		ID:      entity.Id,
		Primary: entity.Primary,
	}
	if len(entity.JobTitle) > 0 {
		contactRole.JobTitle = utils.StringPtr(entity.JobTitle)
	}
	return &contactRole
}

func MapEntitiesToContactRoles(entities *entity.ContactRoleEntities) []*model.ContactRole {
	var contactRoles []*model.ContactRole
	for _, contactRoleEntity := range *entities {
		contactRoles = append(contactRoles, MapEntityToContactRole(&contactRoleEntity))
	}
	return contactRoles
}
