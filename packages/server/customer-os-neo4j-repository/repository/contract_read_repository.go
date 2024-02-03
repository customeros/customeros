package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type ContractReadRepository interface {
	GetContractById(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetContractByServiceLineItemId(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error)
	GetContractByOpportunityId(ctx context.Context, tenant string, opportunityId string) (*dbtype.Node, error)
	GetContractsForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	TenantsHasAtLeastOneContract(ctx context.Context, tenant string) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
	GetContractsToGenerateCycleInvoices(ctx context.Context, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error)
	GetContractsToGenerateOffCycleInvoices(ctx context.Context, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

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

func (r *contractReadRepository) TenantsHasAtLeastOneContract(ctx context.Context, tenant string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.TenantsHasAtLeastOneContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

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

func (r *contractReadRepository) GetContractsToGenerateCycleInvoices(ctx context.Context, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsToGenerateCycleInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Object("referenceTime", referenceTime))

	cypher := `MATCH (ts:TenantSettings)<-[:HAS_SETTINGS]-(t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(:ServiceLineItem)
			WHERE 
				ts.invoicingEnabled = true AND
				(o.hide = false OR o.hide IS NULL) AND
				(c.currency <> "" OR ts.defaultCurrency <> "") AND
				c.organizationLegalName IS NOT NULL AND 
				c.organizationLegalName <> "" AND
				c.invoiceEmail IS NOT NULL AND
				c.invoiceEmail <> "" AND
				(c.billingCycle IN $validBillingCycles) AND
				(c.nextInvoiceDate IS NULL OR date(c.nextInvoiceDate) <= date($referenceTime)) AND
				(date(c.invoicingStartDate) <= date($referenceTime)) AND
				(c.endedAt IS NULL OR date(c.endedAt) > date(coalesce(c.nextInvoiceDate, c.invoicingStartDate))) AND
				(c.techInvoicingStartedAt IS NULL OR c.techInvoicingStartedAt + duration({hours: 1}) < $referenceTime)
			RETURN distinct(c), t.name limit 100`
	params := map[string]any{
		"referenceTime": referenceTime,
		"validBillingCycles": []string{
			neo4jenum.BillingCycleMonthlyBilling.String(), neo4jenum.BillingCycleQuarterlyBilling.String(), neo4jenum.BillingCycleAnnuallyBilling.String(),
		},
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

func (r *contractReadRepository) GetContractsToGenerateOffCycleInvoices(ctx context.Context, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractReadRepository.GetContractsToGenerateOffCycleInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Object("referenceTime", referenceTime))

	cypher := `MATCH (ts:TenantSettings)<-[:HAS_SETTINGS]-(t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract)-[:HAS_SERVICE]->(sli:ServiceLineItem)
			WHERE 
				ts.invoicingEnabled = true AND
				(o.hide = false OR o.hide IS NULL) AND
				(c.currency <> "" OR ts.defaultCurrency <> "") AND
				c.organizationLegalName IS NOT NULL AND 
				c.organizationLegalName <> "" AND
				c.invoiceEmail IS NOT NULL AND
				c.invoiceEmail <> "" AND
				(c.billingCycle IN $validBillingCycles) AND
				c.nextInvoiceDate IS NOT NULL AND
				(c.endedAt IS NULL OR date(c.endedAt) > date($referenceTime)) AND
				NOT EXISTS((sli)<-[:INVOICED]-(:InvoiceLine)) AND
				date(sli.startedAt) + duration({days: 1}) < date(c.nextInvoiceDate) AND
				date(sli.startedAt) < date($referenceTime) AND
				(c.techOffCycleInvoicingStartedAt IS NULL OR date(c.techOffCycleInvoicingStartedAt) < date($referenceTime))
			WITH c, sli, t 
				OPTIONAL MATCH (c)-[:HAS_SERVICE]->(invoicedSli:ServiceLineItem)<-[:INVOICED]-(il:InvoiceLine)
					WHERE EXISTS((invoicedSli)<-[:INVOICED]-(il))
			WITH c, sli, t, invoicedSli
				ORDER BY c, sli, invoicedSli.startedAt DESC
			WITH c, sli, t, head(collect(invoicedSli)) as lastInvoicedSli
				WHERE lastInvoicedSli IS NULL OR date(lastInvoicedSli.startedAt) < date(sli.startedAt) 
			RETURN distinct(c), t.name limit 100`
	params := map[string]any{
		"referenceTime": referenceTime,
		"validBillingCycles": []string{
			neo4jenum.BillingCycleMonthlyBilling.String(), neo4jenum.BillingCycleQuarterlyBilling.String(), neo4jenum.BillingCycleAnnuallyBilling.String(),
		},
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
