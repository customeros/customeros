package neo4j_repository

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/test"
)

func OrganizationWasInserted(ctx context.Context, organizationId string) {
	test.CreateOrganization(ctx, driver, tenantName, neo4j_entity.OrganizationEntity{
		ID: organizationId,
	})
}
