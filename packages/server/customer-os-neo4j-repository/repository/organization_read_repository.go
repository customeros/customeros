package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type OrganizationReadRepository interface {
	GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error)
	GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOrganizationsForInvoicing(ctx context.Context, invoiceDateTime time.Time) ([]map[string]any, error)
}

type organizationReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationReadRepository(driver *neo4j.DriverWithContext, database string) OrganizationReadRepository {
	return &organizationReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *organizationReadRepository) GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$id}) RETURN org`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationIdsConnectedToInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}),
				(t:Tenant {name:$tenant})
				CALL {
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:REPORTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				UNION 
					WITH ie, t 
					MATCH (ie)-[:PART_OF]->(is:Issue)-[:SUBMITTED_BY]->(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
					RETURN org.id as orgId
				}
				RETURN distinct orgId`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), err
}

func (r *organizationReadRepository) GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByOpportunityId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId))

	cypher := `MATCH (op:Opportunity {id:$id})
				MATCH (t:Tenant {name:$tenant})
				OPTIONAL MATCH (op)<-[:HAS_OPPORTUNITY]-(:Contract)<-[:HAS_CONTRACT]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				OPTIONAL MATCH (op)<-[:HAS_OPPORTUNITY]-(directOrg:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			WITH COALESCE(org, directOrg) as organization 
			WHERE organization IS NOT NULL RETURN organization`
	params := map[string]any{
		"tenant": tenant,
		"id":     opportunityId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByContractId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_CONTRACT]->(c:Contract {id:$id})
			RETURN org limit 1`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	if len(records) == 0 {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	} else {
		span.LogFields(log.Bool("result.found", true))
		return records[0], nil
	}
}

func (r *organizationReadRepository) GetOrganizationsForInvoicing(ctx context.Context, invoiceDateTime time.Time) ([]map[string]any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForInvoicing")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Object("invoiceDateTime", invoiceDateTime))

	cypher := fmt.Sprintf(`
			WITH datetime({year: $invoiceDateTime.year, month: $invoiceDateTime.month, day: $invoiceDateTime.day, hour: $invoiceDateTime.hour, minute: $invoiceDateTime.minute, second: $invoiceDateTime.second, nanosecond: $invoiceDateTime.nanosecond}) as invoiceDateTime
			MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			WHERE 
			  o.hide = false AND o.isCustomer = true AND (o.invoicingActive is null or o.invoicingActive = false) AND sli.startedAt < invoiceDateTime AND
			  ((
				(sli.billed = 'MONTHLY' AND (
				  invoiceDateTime.day = sli.startedAt.day
				  OR
				  (
					sli.startedAt.month = 2 AND sli.startedAt.day = 29 AND 
					date({year:invoiceDateTime.year, month:2, day:28}) <= date(invoiceDateTime) AND invoiceDateTime.day = 28
				  )
				))
				OR
				(sli.billed = 'QUARTERLY' AND (
				  (
					date(sli.startedAt).day = date(invoiceDateTime).day AND
					((date(invoiceDateTime).month - date(sli.startedAt).month) %s = 0 OR
					(date(sli.startedAt).month = 2 AND date(invoiceDateTime).month = 2 AND date(sli.startedAt).day = 29 AND date(invoiceDateTime).day = 28))
				  )
				  OR
				  (
					date(sli.startedAt).month = 2 AND date(sli.startedAt).day = 29 AND
					date(invoiceDateTime).month = 2 AND date(invoiceDateTime).day = 28 AND
					date({year: invoiceDateTime.year, month: 2, day: 28}) <= date(invoiceDateTime)
				  )
				))
				OR
				(sli.billed = 'ANNUALLY' AND (
				  date(sli.startedAt).month = date(invoiceDateTime).month AND
				  date(sli.startedAt).day = date(invoiceDateTime).day AND
				  NOT (
					date(sli.startedAt).month = 2 AND date(sli.startedAt).day = 29 AND
					date(invoiceDateTime).month = 2 AND date(invoiceDateTime).day = 28 AND
					date({year: invoiceDateTime.year, month: 2, day: 28}) <= date(invoiceDateTime)
				  )
				))
			  ) OR (datetime({year: sli.startedAt.year, month: sli.startedAt.month, day: sli.startedAt.day}) + duration({days: 1})) <= invoiceDateTime)
			WITH DISTINCT(o.id) as organizationId, t.name as tenant limit 100
			RETURN tenant, organizationId
			`, "% 3")
	params := map[string]any{
		"invoiceDateTime": invoiceDateTime,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}

	response := make([]map[string]any, 0)

	for _, v := range result.([]*neo4j.Record) {
		response = append(response, v.AsMap())
	}
	return response, nil
}
