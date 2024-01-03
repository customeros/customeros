package test

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
	})
}

func CreateMasterPlan(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, masterPlan entity.MasterPlanEntity) string {
	masterPlanId := utils.NewUUIDIfEmpty(masterPlan.Id)

	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$id})
				SET mp:MasterPlan_%s,
					mp.name=$name`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"id":        masterPlanId,
		"name":      masterPlan.Name,
		"createdAt": masterPlan.CreatedAt,
	})
	return masterPlanId
}
