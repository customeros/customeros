package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OrganizationPlanUpdateFields struct {
	Name                string
	Retired             bool
	StatusDetails       entity.OrganizationPlanStatusDetails
	UpdateName          bool
	UpdateRetired       bool
	UpdateStatusDetails bool
}

type OrganizationPlanMilestoneUpdateFields struct {
	Name                string
	Order               int64
	DurationHours       int64
	DueDate             time.Time
	Items               []entity.OrganizationPlanMilestoneItem
	StatusDetails       entity.OrganizationPlanMilestoneStatusDetails
	Optional            bool
	Retired             bool
	Adhoc               bool
	UpdateName          bool
	UpdateOrder         bool
	UpdateItems         bool
	UpdateOptional      bool
	UpdateRetired       bool
	UpdateStatusDetails bool
	UpdateDueDate       bool
	UpdateAdhoc         bool
}

type OrganizationPlanWriteRepository interface {
	Create(ctx context.Context, tenant, masterPlanId, organizationPlanId, name, source, appSource string, createdAt time.Time, statusDetails entity.OrganizationPlanStatusDetails) error
	Update(ctx context.Context, tenant, organizationPlanId string, data OrganizationPlanUpdateFields) error
	CreateMilestone(ctx context.Context, tenant, organizationPlanId, milestoneId, name, source, appSource string, order int64, items []entity.OrganizationPlanMilestoneItem, optional, adhoc bool, createdAt, dueDate time.Time, statusDetails entity.OrganizationPlanMilestoneStatusDetails) error
	CreateBulkMilestones(ctx context.Context, tenant, organizationPlanId, source, appSource string, milestones []entity.OrganizationPlanMilestoneEntity, createdAt time.Time) error
	UpdateMilestone(ctx context.Context, tenant, organizationPlanId, milestoneId string, data OrganizationPlanMilestoneUpdateFields) error
	LinkWithOrganization(ctx context.Context, tenant, organizationPlanId, organizationId string, createdAt time.Time) error
	LinkWithMasterPlan(ctx context.Context, tenant, organizationPlanId, masterPlanId string, createdAt time.Time) error
}

type organizationPlanWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationPlanWriteRepository(driver *neo4j.DriverWithContext, database string) OrganizationPlanWriteRepository {
	return &organizationPlanWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationPlanWriteRepository) Create(ctx context.Context, tenant, masterPlanId, organizationPlanId, name, source, appSource string, createdAt time.Time, statusDetails entity.OrganizationPlanStatusDetails) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							ON CREATE SET 
								op:OrganizationPlan_%s,
								op.createdAt=$createdAt,
								op.updatedAt=datetime(),
								op.source=$source,
								op.sourceOfTruth=$sourceOfTruth,
								op.appSource=$appSource,
								op.name=$name,
								op.status=$status,
								op.statusComments=$statusComments,
								op.statusUpdatedAt=$statusUpdatedAt,
								op.masterPlanId=$masterPlanId
							`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"createdAt":          createdAt,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
		"name":               name,
		"status":             statusDetails.Status,
		"statusComments":     statusDetails.Comments,
		"statusUpdatedAt":    statusDetails.UpdatedAt,
		"masterPlanId":       masterPlanId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationPlanWriteRepository) CreateMilestone(ctx context.Context, tenant, organizationPlanId, milestoneId, name, source, appSource string, order int64, items []entity.OrganizationPlanMilestoneItem, optional, adhoc bool, createdAt, dueDate time.Time, statusDetails entity.OrganizationPlanMilestoneStatusDetails) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.CreateMilestone")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, milestoneId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							MERGE (op)-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:$milestoneId})
							ON CREATE SET 
								m:OrganizationPlanMilestone_%s,
								m.createdAt=$createdAt,
								m.updatedAt=datetime(),
								m.source=$source,
								m.sourceOfTruth=$sourceOfTruth,
								m.appSource=$appSource,
								m.name=$name,
								m.order=$order,
								m.optional=$optional,
								m.items=$items,
								m.status=$status,
								m.statusComments=$statusComments,
								m.statusUpdatedAt=$statusUpdatedAt,
								m.dueDate=$dueDate,
								m.adhoc=$adhoc
							`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"milestoneId":        milestoneId,
		"createdAt":          createdAt,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
		"name":               name,
		"order":              order,
		"optional":           optional,
		"items":              mapMilestoneItemsToNeo4jProperties(items),
		"status":             statusDetails.Status,
		"statusComments":     statusDetails.Comments,
		"statusUpdatedAt":    statusDetails.UpdatedAt,
		"dueDate":            dueDate,
		"adhoc":              adhoc,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationPlanWriteRepository) CreateBulkMilestones(ctx context.Context, tenant, organizationPlanId, source, appSource string, milestones []entity.OrganizationPlanMilestoneEntity, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.CreateBulkMilestones")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "milestones", milestones)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							UNWIND $milestones as milestone
							MERGE (op)-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:milestone.id})
							ON CREATE SET
								m:OrganizationPlanMilestone_%s,
								m.createdAt=$createdAt,
								m.updatedAt=datetime(),
								m.source=$source,
								m.sourceOfTruth=$sourceOfTruth,
								m.appSource=$appSource,
								m.name=milestone.name,
								m.order=milestone.order,
								m.dueDate=milestone.dueDate,
								m.optional=milestone.optional,
								m.items=milestone.items,
								m.status=milestone.status,
								m.statusComments=milestone.statusComments,
								m.statusUpdatedAt=milestone.statusUpdatedAt,
								m.adhoc=milestone.adhoc
							`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"milestones":         mapMilestoneEntitiesToNeo4jProperties(milestones),
		"createdAt":          createdAt,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		err = errors.Wrap(err, "failed to create bulk milestones")
	}
	return err
}

func (r *organizationPlanWriteRepository) Update(ctx context.Context, tenant, organizationPlanId string, data OrganizationPlanUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							SET op.updatedAt=datetime()`
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
	}
	if data.UpdateName {
		cypher += ", op.name=$name"
		params["name"] = data.Name
	}
	if data.UpdateRetired {
		cypher += ", op.retired=$retired"
		params["retired"] = data.Retired
	}
	if data.UpdateStatusDetails {
		cypher += ", op.status=$status, op.statusComments=$statusComments, op.statusUpdatedAt=$statusUpdatedAt"
		params["status"] = data.StatusDetails.Status
		params["statusComments"] = data.StatusDetails.Comments
		params["statusUpdatedAt"] = data.StatusDetails.UpdatedAt
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationPlanWriteRepository) UpdateMilestone(ctx context.Context, tenant, organizationPlanId, milestoneId string, data OrganizationPlanMilestoneUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.UpdateMilestone")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, milestoneId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId})-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:$milestoneId}) 
							SET m.updatedAt=datetime()`
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"milestoneId":        milestoneId,
	}
	if data.UpdateName {
		cypher += ", m.name=$name"
		params["name"] = data.Name
	}
	if data.UpdateOrder {
		cypher += ", m.order=$order"
		params["order"] = data.Order
	}
	if data.UpdateItems {
		cypher += ", m.items=$items"
		params["items"] = mapMilestoneItemsToNeo4jProperties(data.Items)
	}
	if data.UpdateOptional {
		cypher += ", m.optional=$optional"
		params["optional"] = data.Optional
	}
	if data.UpdateRetired {
		cypher += ", m.retired=$retired"
		params["retired"] = data.Retired
	}
	if data.UpdateStatusDetails {
		cypher += ", m.status=$status, m.statusComments=$statusComments, m.statusUpdatedAt=$statusUpdatedAt"
		params["status"] = data.StatusDetails.Status
		params["statusComments"] = data.StatusDetails.Comments
		params["statusUpdatedAt"] = data.StatusDetails.UpdatedAt
	}
	if data.UpdateDueDate {
		cypher += ", m.dueDate=$dueDate"
		params["dueDate"] = data.DueDate
	}
	if data.UpdateAdhoc {
		cypher += ", m.adhoc=$adhoc"
		params["adhoc"] = data.Adhoc
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationPlanWriteRepository) LinkWithOrganization(ctx context.Context, tenant, organizationPlanId, organizationId string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}) 
							MERGE (op)-[:ORGANIZATION_PLAN_BELONGS_TO_ORGANIZATION]->(o)`
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"organizationId":     organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *organizationPlanWriteRepository) LinkWithMasterPlan(ctx context.Context, tenant, organizationPlanId, masterPlanId string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanWriteRepository.LinkWithMasterPlan")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$organizationPlanId}) 
							MATCH (t)<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(m:MasterPlan {id:$masterPlanId}) 
							MERGE (op)-[:ORGANIZATION_PLAN_BELONGS_TO_MASTER_PLAN]->(m)`
	params := map[string]any{
		"tenant":             tenant,
		"organizationPlanId": organizationPlanId,
		"masterPlanId":       masterPlanId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

// func mapPlanEntityToNeo4jProperties(entity entity.OrganizationPlanEntity) map[string]any {
// 	return map[string]any{
// 		"id":              entity.Id,
// 		"createdAt":       entity.CreatedAt,
// 		"updatedAt":       entity.UpdatedAt,
// 		"name":            entity.Name,
// 		"retired":         entity.Retired,
// 		"source":          entity.Source,
// 		"sourceOfTruth":   entity.SourceOfTruth,
// 		"appSource":       entity.AppSource,
// 		"status":          entity.StatusDetails.Status,
// 		"statusComments":  entity.StatusDetails.Comments,
// 		"statusUpdatedAt": entity.StatusDetails.UpdatedAt,
// 	}
// }

func mapMilestoneEntityToNeo4jProperties(entity entity.OrganizationPlanMilestoneEntity) map[string]any {
	return map[string]any{
		"id":              entity.Id,
		"createdAt":       entity.CreatedAt,
		"updatedAt":       entity.UpdatedAt,
		"name":            entity.Name,
		"order":           entity.Order,
		"dueDate":         entity.DueDate,
		"optional":        entity.Optional,
		"retired":         entity.Retired,
		"source":          entity.Source,
		"sourceOfTruth":   entity.SourceOfTruth,
		"appSource":       entity.AppSource,
		"status":          entity.StatusDetails.Status,
		"statusComments":  entity.StatusDetails.Comments,
		"statusUpdatedAt": entity.StatusDetails.UpdatedAt,
		"items":           mapMilestoneItemsToNeo4jProperties(entity.Items),
		"adhoc":           entity.Adhoc,
	}
}

func mapMilestoneItemToNeo4jProperties(item entity.OrganizationPlanMilestoneItem) string {
	ji, _ := json.Marshal(item)
	return string(ji[:]) // fmt.Sprintf(`{"text":%s,"status":%s,"updatedAt":%s,"uuid":%s}`, item.Text, item.Status, item.UpdatedAt, item.Uuid)
}

func mapMilestoneItemsToNeo4jProperties(items []entity.OrganizationPlanMilestoneItem) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = mapMilestoneItemToNeo4jProperties(item)
	}
	return result
}

func mapMilestoneEntitiesToNeo4jProperties(entities []entity.OrganizationPlanMilestoneEntity) []map[string]any {
	result := make([]map[string]any, len(entities))
	for i, entity := range entities {
		result[i] = mapMilestoneEntityToNeo4jProperties(entity)
	}
	return result
}
