package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapTenantInputToEntity(input model.TenantInput) entity.TenantEntity {
	tenantEntity := entity.TenantEntity{
		Name:      input.Name,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: utils.IfNotNilString(input.AppSource),
	}
	if len(tenantEntity.AppSource) == 0 {
		tenantEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return tenantEntity
}

func MapEntitiesToTenantNames(entities *entity.TenantEntities) []string {
	var tenants []string
	for _, tenantEntity := range *entities {
		tenants = append(tenants, tenantEntity.Name)
	}
	return tenants
}
