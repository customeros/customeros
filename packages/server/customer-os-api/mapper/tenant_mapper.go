package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
)

func MapEntitiesToTenantNames(entities *entity.TenantEntities) []string {
	var tenants []string
	for _, tenantEntity := range *entities {
		tenants = append(tenants, tenantEntity.Name)
	}
	return tenants
}
