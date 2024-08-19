package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type TenantAndContractId struct {
	Tenant     string
	ContractId string
}

type ContractReadRepository interface {
	GetContractById(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetContractByServiceLineItemId(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error)
	GetContractByOpportunityId(ctx context.Context, tenant string, opportunityId string) (*dbtype.Node, error)
	GetContractsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetContractForInvoice(ctx context.Context, tenant string, invoiceId string) (*dbtype.Node, error)
	GetContractsForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error)
	TenantsHasAtLeastOneContract(ctx context.Context, tenant string) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
	GetContractsToGenerateCycleInvoices(ctx context.Context, referenceTime time.Time, delayMinutes, limit int) ([]*utils.DbNodeAndTenant, error)
	GetContractsToGenerateOffCycleInvoices(ctx context.Context, referenceTime time.Time, delayMinutes, limit int) ([]*utils.DbNodeAndTenant, error)
	GetContractsToGenerateNextScheduledInvoices(ctx context.Context, referenceTime time.Time, delayMinutes int) ([]*utils.DbNodeAndTenant, error)
	GetContractsForStatusRenewal(ctx context.Context, referenceTime time.Time, limit int) ([]TenantAndContractId, error)
	GetContractsForRenewalRollout(ctx context.Context, referenceTime time.Time, limit int) ([]TenantAndContractId, error)
	IsContractInvoiced(ctx context.Context, tenant, contractId string) (bool, error)
	GetPaginatedContracts(ctx context.Context, tenant string, skip, limit int) (*utils.DbNodesWithTotalCount, error)
	GetLiveContractsWithoutRenewalOpportunities(ctx context.Context, limit int) ([]TenantAndContractId, error)
}

type contractReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewContractReadRepository(driver *neo4j.DriverWithContext, database string) ContractReadRepository {
	return &contractReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *contractReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *contractReadRepository) GetContractById(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	span.LogFields(log.String("contractId", contractId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *contractReadRepository) GetContractByServiceLineItemId(ctx context.Context, tenant string, serviceLineItemId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractByServiceLineItemId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	cypher := `MATCH (sli:ServiceLineItem {id:$id})<-[:HAS_SERVICE]-(c:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) RETURN c LIMIT 1`
	params := map[string]any{
		"id":     serviceLineItemId,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	span.LogFields(log.Int("result.count", len(records)))
	if len(records) == 0 {
		return nil, nil
	} else {
		return records[0], nil
	}
}

func (r *contractReadRepository) GetContractByOpportunityId(ctx context.Context, tenant string, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractByOpportunityId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("opportunityId", opportunityId))

	cypher := fmt.Sprintf(`MATCH (:Opportunity {id:$id})<-[:HAS_OPPORTUNITY]-(c:Contract:Contract_%s)-[:CONTRACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) RETURN c LIMIT 1`, tenant)
	params := map[string]any{
		"id":     opportunityId,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	span.LogFields(log.Int("result.count", len(records)))
	if len(records) == 0 {
		return nil, nil
	} else {
		return records[0], nil
	}
}

func (r *contractReadRepository) GetContractsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsForOrganizations")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(contract:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN contract, o.id ORDER BY contract.createdAt DESC`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *contractReadRepository) GetContractForInvoice(ctx context.Context, tenant string, invoiceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractForInvoice")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_INVOICE]->(i:Invoice_%s)
			WHERE i.id = $invoiceId
			RETURN c`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	records := result.([]*dbtype.Node)
	span.LogFields(log.Int("result.count", len(records)))
	if len(records) == 0 {
		return nil, nil
	} else {
		return records[0], nil
	}
}

func (r *contractReadRepository) GetContractsForInvoices(ctx context.Context, tenant string, invoiceIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsForInvoices")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_INVOICE]->(i:Invoice_%s)
			WHERE i.id IN $invoiceIds
			RETURN c, i.id`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"invoiceIds": invoiceIds,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *contractReadRepository) TenantsHasAtLeastOneContract(ctx context.Context, tenant string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.TenantsHasAtLeastOneContract")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (c:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(:Tenant{name:$tenant}) RETURN count(c) > 0`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), err
}

func (r *contractReadRepository) CountContracts(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.CountContracts")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (c:Contract)-[:CONTRACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})MATCH (c)-[:ACTIVE_RENEWAL]->(op:Opportunity)
			RETURN count(c)`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	contractsCount := dbRecord.(*db.Record).Values[0].(int64)
	span.LogFields(log.Int64("result - contractsCount", contractsCount))
	return contractsCount, nil
}

func (r *contractReadRepository) GetContractsToGenerateCycleInvoices(ctx context.Context, referenceTime time.Time, delayMinutes, limit int) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsToGenerateCycleInvoices")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Object("referenceTime", referenceTime), log.Int("delayMinutes", delayMinutes), log.Int("limit", limit))

	cypher := `MATCH (ts:TenantSettings)<-[:HAS_SETTINGS]-(t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(:ServiceLineItem)
			WHERE 
				ts.invoicingEnabled = true AND
				(c.invoicingEnabled = true OR c.invoicingEnabled IS NULL) AND
				(o.hide = false OR o.hide IS NULL) AND
				(c.currency <> "" OR ts.baseCurrency <> "" ) AND
				c.organizationLegalName IS NOT NULL AND 
				c.organizationLegalName <> "" AND
				c.invoiceEmail IS NOT NULL AND
				c.invoiceEmail <> "" AND
				c.billingCycleInMonths > 0 AND
				c.status IN $validContractStatuses AND
				(c.nextInvoiceDate IS NULL OR date(c.nextInvoiceDate) <= date($referenceTime)) AND
				(date(c.invoicingStartDate) <= date($referenceTime)) AND
				(c.endedAt IS NULL OR date(c.endedAt) > date(coalesce(c.nextInvoiceDate, c.invoicingStartDate))) AND
				(c.techInvoicingStartedAt IS NULL OR c.techInvoicingStartedAt + duration({minutes: $delayMinutes}) < $referenceTime)
			RETURN distinct(c), t.name limit $limit`
	params := map[string]any{
		"referenceTime":         referenceTime,
		"validContractStatuses": []string{neo4jenum.ContractStatusLive.String()},
		"delayMinutes":          delayMinutes,
		"limit":                 limit,
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
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndTenant))))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *contractReadRepository) GetContractsToGenerateOffCycleInvoices(ctx context.Context, referenceTime time.Time, delayMinutes, limit int) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsToGenerateOffCycleInvoices")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Object("referenceTime", referenceTime), log.Int("delayMinutes", delayMinutes), log.Int("limit", limit))

	cypher := `MATCH (ts:TenantSettings)<-[:HAS_SETTINGS]-(t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			WHERE 
				ts.invoicingEnabled = true AND 
				(ts.invoicingPostpaid = false OR ts.invoicingPostpaid IS NULL) AND 
				(c.invoicingEnabled = true OR c.invoicingEnabled IS NULL) AND
				(o.hide = false OR o.hide IS NULL) AND
				(c.currency <> "" OR ts.baseCurrency <> "") AND
				c.organizationLegalName IS NOT NULL AND 
				c.organizationLegalName <> "" AND
				c.invoiceEmail IS NOT NULL AND
				c.invoiceEmail <> "" AND
				c.billingCycleInMonths > 0 AND
				c.status IN $validContractStatuses AND
				c.nextInvoiceDate IS NOT NULL AND
				(c.endedAt IS NULL OR date(c.endedAt) > date($referenceTime)) AND
				NOT EXISTS((sli)<-[:INVOICED]-(:InvoiceLine)--(:Invoice {dryRun:false})) AND
				date(sli.startedAt) + duration({days: 1}) < date(c.nextInvoiceDate) AND
				date(sli.startedAt) < date($referenceTime) AND
				(sli.isCanceled = false OR sli.isCanceled IS NULL) AND
				(c.techOffCycleInvoicingStartedAt IS NULL OR date(c.techOffCycleInvoicingStartedAt) < date($referenceTime) OR c.techOffCycleInvoicingStartedAt + duration({minutes: $delayMinutes}) < $referenceTime)
			WITH c, sli, t 
				OPTIONAL MATCH (c)-[:HAS_SERVICE]->(invoicedSli:ServiceLineItem)<-[:INVOICED]-(il:InvoiceLine)
					WHERE EXISTS((invoicedSli)<-[:INVOICED]-(il)--(:Invoice {dryRun:false}))
			WITH c, sli, t, invoicedSli
				ORDER BY c, sli, invoicedSli.startedAt DESC
			WITH c, sli, t, head(collect(invoicedSli)) as lastInvoicedSli
				WHERE lastInvoicedSli IS NULL OR date(lastInvoicedSli.startedAt) < date(sli.startedAt) 
			RETURN distinct(c), t.name limit $limit`
	params := map[string]any{
		"referenceTime":         referenceTime,
		"validContractStatuses": []string{neo4jenum.ContractStatusLive.String()},
		"delayMinutes":          delayMinutes,
		"limit":                 limit,
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
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndTenant))))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *contractReadRepository) GetContractsToGenerateNextScheduledInvoices(ctx context.Context, referenceTime time.Time, delayMinutes int) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsToGenerateNextScheduledInvoices")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Object("referenceTime", referenceTime), log.Int("delayMinutes", delayMinutes))

	cypher := `MATCH (ts:TenantSettings)<-[:HAS_SETTINGS]-(t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			OPTIONAL MATCH (c)-[:HAS_INVOICE]->(i:Invoice {dryRun: true, preview: true})
			WITH c, t, ts, o, i
			WHERE 
				(i IS NULL OR i.createdAt < c.updatedAt OR i.createdAt < sli.updatedAt) AND
				ts.invoicingEnabled = true AND
				(c.invoicingEnabled = true OR c.invoicingEnabled IS NULL) AND
				(o.hide = false OR o.hide IS NULL) AND
				(c.currency <> "" OR ts.baseCurrency <> "" ) AND
				c.organizationLegalName IS NOT NULL AND 
				c.organizationLegalName <> "" AND
				c.invoiceEmail IS NOT NULL AND
				c.invoiceEmail <> "" AND
				c.billingCycleInMonths > 0 AND
				c.status IN $validContractStatuses AND
				(NOT c.invoicingStartDate IS NULL OR NOT c.nextInvoiceDate IS NULL) AND
				(c.endedAt IS NULL OR date(c.endedAt) > date(coalesce(c.nextInvoiceDate, c.invoicingStartDate))) AND
				(c.techNextPreviewInvoiceRequestedAt IS NULL OR c.techNextPreviewInvoiceRequestedAt + duration({minutes: $delayMinutes}) < $referenceTime)
			RETURN distinct(c), t.name limit 100`
	params := map[string]any{
		"referenceTime": referenceTime,
		"validContractStatuses": []string{neo4jenum.ContractStatusLive.String(),
			neo4jenum.ContractStatusOutOfContract.String(),
			neo4jenum.ContractStatusScheduled.String(),
			neo4jenum.ContractStatusDraft.String()},
		"delayMinutes": delayMinutes,
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
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndTenant))))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *contractReadRepository) GetContractsForStatusRenewal(ctx context.Context, referenceTime time.Time, limit int) ([]TenantAndContractId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetContractsForStatusRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Object("referenceTime", referenceTime), log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)<-[:HAS_CONTRACT]-(:Organization {hide:false})
				WHERE c.techStatusRenewalRequestedAt IS NULL OR c.techStatusRenewalRequestedAt + duration({hours: 2}) < $referenceTime
				OPTIONAL MATCH (c)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WITH t, c, op.renewedAt as renewedAt
				WHERE (c.status <> $endedStatus AND c.endedAt < $referenceTime) OR
						((c.endedAt IS NULL OR c.endedAt > $referenceTime) AND 
						(
							(c.status in [$scheduledStatus, $endedStatus] AND date(c.serviceStartedAt) <= date($referenceTime)) OR 
							(c.status = $outOfContractStatus AND (c.autoRenew = true OR renewedAt > $referenceTime)) OR
							(c.status = $outOfContractStatus AND renewedAt > $referenceTime) OR
							(c.status = $liveStatus AND c.autoRenew = false AND renewedAt < $referenceTime) OR 
							(c.status = $draftStatus AND c.approved = true)
						))
				RETURN DISTINCT t.name, c.id LIMIT $limit`
	params := map[string]any{
		"referenceTime":       referenceTime,
		"endedStatus":         neo4jenum.ContractStatusEnded,
		"liveStatus":          neo4jenum.ContractStatusLive,
		"scheduledStatus":     neo4jenum.ContractStatusScheduled,
		"outOfContractStatus": neo4jenum.ContractStatusOutOfContract,
		"draftStatus":         neo4jenum.ContractStatusDraft,
		"limit":               limit,
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
	output := make([]TenantAndContractId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContractId{
				Tenant:     v.Values[0].(string),
				ContractId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contractReadRepository) GetContractsForRenewalRollout(ctx context.Context, referenceTime time.Time, limit int) ([]TenantAndContractId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetContractsForStatusRenewal")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Object("referenceTime", referenceTime), log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract),
				(c)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE
					(c.techRolloutRenewalRequestedAt IS NULL OR c.techRolloutRenewalRequestedAt + duration({hours: 2}) < $referenceTime) AND
					date(op.renewedAt) <= date($referenceTime) AND
					(c.autoRenew = true OR op.renewalApproved = true) AND
					c.status IN [$liveStatus, $outOfContractStatus, $scheduledStatus]
				RETURN t.name, c.id LIMIT $limit`
	params := map[string]any{
		"referenceTime":       referenceTime,
		"liveStatus":          neo4jenum.ContractStatusLive.String(),
		"outOfContractStatus": neo4jenum.ContractStatusOutOfContract.String(),
		"scheduledStatus":     neo4jenum.ContractStatusScheduled.String(),
		"limit":               limit,
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
	output := make([]TenantAndContractId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContractId{
				Tenant:     v.Values[0].(string),
				ContractId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *contractReadRepository) IsContractInvoiced(ctx context.Context, tenant, contractId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.IsContractInvoiced")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice {dryRun:false})
			RETURN count(i) > 0`
	params := map[string]any{
		"contractId": contractId,
		"tenant":     tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), err
}

func (r *contractReadRepository) GetPaginatedContracts(ctx context.Context, tenant string, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.IsContractInvoiced")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	countCypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)<-[:HAS_CONTRACT]-(:Organization {hide:false}) RETURN count(c) as count`
	countParams := map[string]any{
		"tenant": tenant,
	}

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)<-[:HAS_CONTRACT]-(:Organization {hide:false}) RETURN c SKIP $skip LIMIT $limit`
	params := map[string]any{
		"tenant": tenant,
		"skip":   skip,
		"limit":  limit,
	}
	span.LogFields(log.String("countCypher", countCypher))
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)
	tracing.LogObjectAsJson(span, "countParams", countParams)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, countCypher, countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, cypher, params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *contractReadRepository) GetLiveContractsWithoutRenewalOpportunities(ctx context.Context, limit int) ([]TenantAndContractId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetLiveContractsWithoutRenewalOpportunities")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	span.LogFields(log.Int("limit", limit))

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)
				WHERE c.status = $liveStatus AND c.lengthInMonths > 0  AND NOT (c)-[:ACTIVE_RENEWAL]->(:RenewalOpportunity)
				RETURN c, t.name LIMIT $limit`
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
	output := make([]TenantAndContractId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndContractId{
				Tenant:     v.Values[0].(string),
				ContractId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
