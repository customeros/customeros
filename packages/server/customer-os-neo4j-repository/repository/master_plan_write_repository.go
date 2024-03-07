package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type MasterPlanUpdateFields struct {
	Name          string
	Retired       bool
	UpdatedAt     time.Time
	UpdateName    bool
	UpdateRetired bool
}

type MasterPlanMilestoneUpdateFields struct {
	UpdatedAt           time.Time
	Name                string
	Order               int64
	DurationHours       int64
	Items               []string
	Optional            bool
	Retired             bool
	UpdateName          bool
	UpdateOrder         bool
	UpdateItems         bool
	UpdateOptional      bool
	UpdateRetired       bool
	UpdateDurationHours bool
}

type MasterPlanWriteRepository interface {
	Create(ctx context.Context, tenant, masterPlanId, name, source, appSource string, createdAt time.Time) error
	Update(ctx context.Context, tenant, masterPlanId string, data MasterPlanUpdateFields) error
	CreateMilestone(ctx context.Context, tenant, masterPlanId, milestoneId, name, source, appSource string, order, durationHours int64, items []string, optional bool, createdAt time.Time) error
	UpdateMilestone(ctx context.Context, tenant, masterPlanId, milestoneId string, data MasterPlanMilestoneUpdateFields) error
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

func (r *masterPlanWriteRepository) CreateMilestone(ctx context.Context, tenant, masterPlanId, milestoneId, name, source, appSource string, order, durationHours int64, items []string, optional bool, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanWriteRepository.CreateMilestone")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, milestoneId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$masterPlanId}) 
							MERGE (mp)-[:HAS_MILESTONE]->(m:MasterPlanMilestone {id:$milestoneId})
							ON CREATE SET 
								m:MasterPlanMilestone_%s,
								m.createdAt=$createdAt,
								m.updatedAt=$updatedAt,
								m.source=$source,
								m.sourceOfTruth=$sourceOfTruth,
								m.appSource=$appSource,
								m.name=$name,
								m.order=$order,
								m.durationHours=$durationHours,
								m.optional=$optional,
								m.items=$items
							`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"masterPlanId":  masterPlanId,
		"milestoneId":   milestoneId,
		"createdAt":     createdAt,
		"updatedAt":     createdAt,
		"source":        source,
		"sourceOfTruth": source,
		"appSource":     appSource,
		"name":          name,
		"order":         order,
		"durationHours": durationHours,
		"optional":      optional,
		"items":         items,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *masterPlanWriteRepository) Update(ctx context.Context, tenant, masterPlanId string, data MasterPlanUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$masterPlanId}) 
							SET mp.updatedAt=$updatedAt`
	params := map[string]any{
		"tenant":       tenant,
		"masterPlanId": masterPlanId,
		"updatedAt":    data.UpdatedAt,
	}
	if data.UpdateName {
		cypher += ", mp.name=$name"
		params["name"] = data.Name
	}
	if data.UpdateRetired {
		cypher += ", mp.retired=$retired"
		params["retired"] = data.Retired
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *masterPlanWriteRepository) UpdateMilestone(ctx context.Context, tenant, masterPlanId, milestoneId string, data MasterPlanMilestoneUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanWriteRepository.UpdateMilestone")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, milestoneId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$masterPlanId})-[:HAS_MILESTONE]->(m:MasterPlanMilestone {id:$milestoneId}) 
							SET m.updatedAt=$updatedAt`
	params := map[string]any{
		"tenant":       tenant,
		"masterPlanId": masterPlanId,
		"milestoneId":  milestoneId,
		"updatedAt":    data.UpdatedAt,
	}
	if data.UpdateName {
		cypher += ", m.name=$name"
		params["name"] = data.Name
	}
	if data.UpdateOrder {
		cypher += ", m.order=$order"
		params["order"] = data.Order
	}
	if data.UpdateDurationHours {
		cypher += ", m.durationHours=$durationHours"
		params["durationHours"] = data.DurationHours
	}
	if data.UpdateItems {
		cypher += ", m.items=$items"
		params["items"] = data.Items
	}
	if data.UpdateOptional {
		cypher += ", m.optional=$optional"
		params["optional"] = data.Optional
	}
	if data.UpdateRetired {
		cypher += ", m.retired=$retired"
		params["retired"] = data.Retired
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
