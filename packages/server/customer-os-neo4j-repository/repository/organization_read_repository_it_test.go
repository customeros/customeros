package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"testing"
)

func TestOrganizationReadRepository_GetOrganizationsForInvoicing(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)

	organization1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, entity.OrganizationEntity{Name: "org 1"})
	neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, organization1Id, entity.ContractEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{
		neo4jutil.NodeLabelOrganization: 1,
		neo4jutil.NodeLabelContract:     1})
}
