package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"strings"
)

func MapTenantInputToEntity(input model.TenantInput) neo4jentity.TenantEntity {
	tenantEntity := neo4jentity.TenantEntity{
		Name:      strings.TrimSpace(input.Name),
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: utils.IfNotNilString(input.AppSource),
	}
	if len(tenantEntity.AppSource) == 0 {
		tenantEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return tenantEntity
}
