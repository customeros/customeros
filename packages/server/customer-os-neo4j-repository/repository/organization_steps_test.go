package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
)

func OrganizationWasInserted(ctx context.Context, organizationId string) {
	test.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		ID: organizationId,
	})
}
