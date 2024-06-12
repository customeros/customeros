package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantAndOrganizationId struct {
	Tenant         string
	OrganizationId string
}

type TenantAndOrganizationIdExtended struct {
	Tenant         string
	OrganizationId string
	Param1         string
}

type OrganizationReadRepository interface {
	GetOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetOrganizationIdsConnectedToInteractionEvent(ctx context.Context, tenant, interactionEventId string) ([]string, error)
	GetOrganizationByOpportunityId(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	GetOrganizationByContractId(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOrganizationByInvoiceId(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error)
	GetOrganizationByCustomerOsId(ctx context.Context, tenant, customerOsId string) (*dbtype.Node, error)
	GetOrganizationByReferenceId(ctx context.Context, tenant, referenceId string) (*dbtype.Node, error)
	GetOrganizationWithDomain(ctx context.Context, tenant, domain string) (*dbtype.Node, error)
	GetAllForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error)
	GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error)
	GetOrganizationsForEnrich(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error)
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

func (r *organizationReadRepository) GetOrganizationByInvoiceId(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByInvoiceId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("invoiceId", invoiceId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(inv:Invoice {id:$invoiceId})<-[:HAS_INVOICE]-(c:Contract)<-[:HAS_CONTRACT]-(org:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
			RETURN org`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
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

func (r *organizationReadRepository) GetOrganizationByCustomerOsId(ctx context.Context, tenant, customerOsId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByInvoiceId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("customerOsId", customerOsId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {customerOsId:$customerOsId})
			RETURN org`
	params := map[string]any{
		"tenant":       tenant,
		"customerOsId": customerOsId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationByReferenceId(ctx context.Context, tenant, referenceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationByReferenceId")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("referenceId", referenceId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {referenceId:$referenceId}) RETURN org`
	params := map[string]any{
		"tenant":      tenant,
		"referenceId": referenceId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))
	return result.(*dbtype.Node), nil
}

func (r *organizationReadRepository) GetOrganizationWithDomain(ctx context.Context, tenant, domain string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationWithDomain")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("domain", domain))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain{domain:$domain}) RETURN o`

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, map[string]any{
			"tenant": tenant,
			"domain": domain,
		}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", true))

	return result.(*dbtype.Node), err
}

func (r *organizationReadRepository) GetAllForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(:Contract)<-[:HAS_CONTRACT]-(o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
				WHERE i.id IN $invoiceIds
				RETURN o, i.id`
	params := map[string]any{
		"tenant":     tenant,
		"invoiceIds": invoiceIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetAllForSlackChannels(ctx context.Context, tenant string, slackChannelIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForSlackChannels")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.Object("slackChannelIds", slackChannelIds))

	cypher := `MATCH (t:Tenant {name:$tenant})-[:ORGANIZATION_BELONGS_TO_TENANT]->(o:Organization)
				WHERE o.slackChannelId IN $slackChannelIds
				RETURN o, i.id`
	params := map[string]any{
		"tenant":          tenant,
		"slackChannelIds": slackChannelIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetAllForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetAllForOpportunities")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	cypher := `MATCH (t:Tenant {name:$tenant})-[:ORGANIZATION_BELONGS_TO_TENANT]->(org:Organization)-[:HAS_OPPORTUNITY]->(op:Opportunity)
				WHERE op.id IN $opportunityIds
				RETURN org, op.id`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *organizationReadRepository) GetOrganizationsForUpdateNextRenewalDate(ctx context.Context, limit int) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_CONTRACT]-(c:Contract)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity) 
				WITH t, org, collect(c) as contracts, collect(op) as ops 
					WHERE ALL(c IN contracts WHERE c.status = $liveStatus) 
				UNWIND ops AS op
				WITH t, org, min(op.renewedAt) AS minOpRenewalDate 
					WHERE date(org.derivedNextRenewalAt) < date(minOpRenewalDate) 
				RETURN t.name, org.id LIMIT $limit`
	params := map[string]any{
		"liveStatus": neo4jenum.ContractStatusLive.String(),
		"limit":      limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsWithWebsiteAndWithoutDomains(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsWithWebsiteAndWithoutDomains")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("limit", limit), log.Int("delayInMinutes", delayInMinutes))

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization) 
				WHERE NOT (org)-[:HAS_DOMAIN]->(:Domain) AND 
						org.website IS NOT NULL AND 
						org.website <> "" AND 
						(org.techDomainCheckedAt IS NULL OR org.techDomainCheckedAt < datetime() - duration({minutes: $delayInMinutes}))
				RETURN t.name, org.id ORDER BY org.createdAt DESC LIMIT $limit`
	params := map[string]any{
		"limit":          limit,
		"delayInMinutes": delayInMinutes,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *organizationReadRepository) GetOrganizationsForEnrich(ctx context.Context, limit, delayInMinutes int) ([]TenantAndOrganizationIdExtended, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationReadRepository.GetOrganizationsForEnrich")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("limit", limit), log.Int("delayInMinutes", delayInMinutes))

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)-[:HAS_DOMAIN]->(d:Domain)
				WHERE org.enrichedAt IS NULL AND
						org.hide = false AND
						(NOT d.enrichedAt IS NULL OR d.enrichRequestedAt IS NULL) AND
						(org.techDomainCheckedAt IS NULL OR org.techDomainCheckedAt < datetime() - duration({minutes: $delayInMinutes}))
				RETURN t.name, org.id, d.domain ORDER BY org.createdAt DESC LIMIT $limit`
	params := map[string]any{
		"limit":          limit,
		"delayInMinutes": delayInMinutes,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndOrganizationIdExtended, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationIdExtended{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
				Param1:         v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
