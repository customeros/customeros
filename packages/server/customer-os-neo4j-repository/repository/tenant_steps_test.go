package repository

import (
	"context"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
)

func TenantWasInserted(ctx context.Context) {
	neo4jtest.CreateTenant(ctx, driver, tenantName)
}
