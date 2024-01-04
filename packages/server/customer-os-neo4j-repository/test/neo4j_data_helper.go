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
	query := `MERGE (t:Tenant {name:$tenant}) ON CREATE SET t.createdAt=$now`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
		"now":    utils.Now(),
	})
}

func CreateMasterPlan(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, masterPlan entity.MasterPlanEntity) string {
	masterPlanId := utils.NewUUIDIfEmpty(masterPlan.Id)

	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$id})
				SET mp:MasterPlan_%s,
					mp.name=$name,
					mp.createdAt=$createdAt,
					mp.source=$source,
					mp.sourceOfTruth=$sourceOfTruth,
					mp.appSource=$appSource
					`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            masterPlanId,
		"name":          masterPlan.Name,
		"createdAt":     masterPlan.CreatedAt,
		"source":        masterPlan.Source,
		"sourceOfTruth": masterPlan.SourceOfTruth,
		"appSource":     masterPlan.AppSource,
	})
	return masterPlanId
}

func CreateMasterPlanMilestone(ctx context.Context, driver *neo4j.DriverWithContext, tenant, masterPlanId string, masterPlanMilestone entity.MasterPlanMilestoneEntity) string {
	masterPlanMilestoneId := utils.NewUUIDIfEmpty(masterPlanMilestone.Id)

	query := fmt.Sprintf(`MATCH (mp:MasterPlan {id: $masterPlanId})
			  MERGE (mp)-[:HAS_MILESTONE]->(m:MasterPlanMilestone {id:$id})
				SET m:MasterPlanMilestone_%s,
					m.name=$name,
					m.createdAt=$createdAt,
					m.order=$order,
					m.durationHours=$durationHours,
					m.optional=$optional,
					m.items=$items`, tenant)

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"masterPlanId":  masterPlanId,
		"id":            masterPlanMilestoneId,
		"name":          masterPlanMilestone.Name,
		"createdAt":     masterPlanMilestone.CreatedAt,
		"order":         masterPlanMilestone.Order,
		"durationHours": masterPlanMilestone.DurationHours,
		"optional":      masterPlanMilestone.Optional,
		"items":         masterPlanMilestone.Items,
	})
	return masterPlanMilestoneId
}
