package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type MasterPlanWriteRepository interface {
	Create(ctx context.Context, tenant, masterPlanId, name, source, appSource string, createdAt time.Time) error
}

type masterPlanWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewMasterPlanWriteRepository(driver *neo4j.DriverWithContext, database string) MasterPlanWriteRepository {
	return &masterPlanWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *masterPlanWriteRepository) Create(ctx context.Context, tenant, masterPlanId, name, source, appSource string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$masterPlanId}) 
							ON CREATE SET 
								mp:MasterPlan_%s,
								mp.createdAt=$createdAt,
								mp.updatedAt=$updatedAt,
								mp.source=$source,
								mp.sourceOfTruth=$sourceOfTruth,
								mp.appSource=$appSource,
								mp.name=$name
							`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"masterPlanId":  masterPlanId,
		"createdAt":     createdAt,
		"updatedAt":     createdAt,
		"source":        source,
		"sourceOfTruth": source,
		"appSource":     appSource,
		"name":          name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
