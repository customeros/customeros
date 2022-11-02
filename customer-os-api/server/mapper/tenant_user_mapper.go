package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapTenantUserInputToEntity(input model.TenantUserInput) *entity.TenantUserEntity {
	tenantUserEntity := entity.TenantUserEntity{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}
	return &tenantUserEntity
}

func MapEntityToTenantUser(tenantUserEntity *entity.TenantUserEntity) *model.TenantUser {
	return &model.TenantUser{
		ID:        tenantUserEntity.Id,
		FirstName: tenantUserEntity.FirstName,
		LastName:  tenantUserEntity.LastName,
		Email:     tenantUserEntity.Email,
		CreatedAt: tenantUserEntity.CreatedAt,
	}
}

func MapEntitiesToTenantUsers(tenantUserEntities *entity.TenantUserEntities) []*model.TenantUser {
	var tenantUsers []*model.TenantUser
	for _, tenantUserEntity := range *tenantUserEntities {
		tenantUsers = append(tenantUsers, MapEntityToTenantUser(&tenantUserEntity))
	}
	return tenantUsers
}
