package repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
)

func OrganizationWasInserted(ctx context.Context, organizationId string) {
	test.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{
		ID: organizationId,
	})
}
