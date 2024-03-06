package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type InvoiceReadRepository interface {
	GetInvoiceById(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error)
	GetPaginatedInvoices(ctx context.Context, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetInvoicesForPayNotifications(ctx context.Context, minutesFromLastUpdate, invoiceNotOlderThanDays int, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error)
	CountNonDryRunInvoicesForContract(ctx context.Context, tenant, contractId string) (int, error)
	GetInvoicesForPaymentLinkRequest(ctx context.Context, minutesFromLastUpdate int, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error)
}

type invoiceReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoiceReadRepository(driver *neo4j.DriverWithContext, database string) InvoiceReadRepository {
	return &invoiceReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoiceReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *invoiceReadRepository) GetPaginatedInvoices(ctx context.Context, tenant, organizationId string, skip, limit int, filter *utils.CypherFilter, sorting *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.GetPaginatedInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.Int("skip", skip))
	span.LogFields(log.Int("limit", limit))
	span.LogFields(log.Object("filter", filter))
	span.LogFields(log.Object("sorting", sorting))

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("i")

		countParams := map[string]any{
			"tenant": tenant,
		}
		queryParams := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}

		if organizationId != "" {
			if filterCypherStr != "" {
				filterCypherStr += " AND "
			} else {
				filterCypherStr += " WHERE "
			}
			filterCypherStr += " o.id=$organizationId"
			filterParams["organizationId"] = organizationId
		}

		utils.MergeMapToMap(filterParams, countParams)

		countCypher := fmt.Sprintf(` MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) 
				 %s 
				 RETURN count(i) as count`, tenant, tenant, tenant, filterCypherStr)

		span.LogFields(log.String("countCypher", countCypher))
		tracing.LogObjectAsJson(span, "countParams", countParams)

		queryResult, err := tx.Run(ctx, countCypher, countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		utils.MergeMapToMap(filterParams, queryParams)

		cypher := fmt.Sprintf(` MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) 
				 %s 
				 RETURN i 
				 %s 
				 SKIP $skip LIMIT $limit`, tenant, tenant, tenant, filterCypherStr, sorting.SortingCypherFragment("i"))

		span.LogFields(log.String("cypher", cypher))
		tracing.LogObjectAsJson(span, "queryParams", queryParams)

		queryResult, err = tx.Run(ctx, cypher,
			queryParams)
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

func (r *invoiceReadRepository) GetInvoiceById(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.GetInvoiceById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) RETURN i`
	params := map[string]any{
		"tenant": tenant,
		"id":     invoiceId,
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

func (r *invoiceReadRepository) GetInvoicesForPayNotifications(ctx context.Context, minutesFromLastUpdate, invoiceNotOlderThanDays int, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.GetInvoicesForPayNotifications")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("minutesFromLastUpdate", minutesFromLastUpdate), log.Int("invoiceNotOlderThanDays", invoiceNotOlderThanDays), log.Object("referenceTime", referenceTime))

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = false AND
				NOT i.status IN $ignoredStatuses AND
				(i.techPayNotificationRequestedAt IS NULL OR i.techPayNotificationRequestedAt + duration({hours: 1}) < $referenceTime) AND
				i.customerEmail IS NOT NULL AND
				i.customerEmail <> '' AND	
				i.techPayInvoiceNotificationSentAt IS NULL AND
				i.createdAt+duration({days: $lookbackWindow}) > $now AND
				(i.updatedAt + duration({minutes: $delay}) < $referenceTime)
			RETURN distinct(i), t.name limit 100`
	params := map[string]any{
		"delay":          minutesFromLastUpdate,
		"lookbackWindow": invoiceNotOlderThanDays,
		"referenceTime":  referenceTime,
		"now":            utils.Now(),
		"ignoredStatuses": []string{
			neo4jenum.InvoiceStatusPaid.String(), neo4jenum.InvoiceStatusDraft.String(), neo4jenum.InvoiceStatusNone.String(),
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

func (r *invoiceReadRepository) CountNonDryRunInvoicesForContract(ctx context.Context, tenant, contractId string) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.CountNonDryRunInvoicesForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)
	span.LogFields(log.String("contractId", contractId))

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice {dryRun:false}) RETURN count(i) as count`
		params := map[string]any{
			"tenant":     tenant,
			"contractId": contractId,
		}
		span.LogFields(log.String("cypher", cypher))
		tracing.LogObjectAsJson(span, "params", params)

		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		return 0, err
	}
	span.LogFields(log.Int64("result.count", count.(int64)))
	return int(count.(int64)), nil
}

func (r *invoiceReadRepository) GetInvoicesForPaymentLinkRequest(ctx context.Context, minutesFromLastUpdate int, referenceTime time.Time) ([]*utils.DbNodeAndTenant, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.GetInvoicesForPaymentLinkRequest")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, "")
	span.LogFields(log.Int("minutesFromLastUpdate", minutesFromLastUpdate), log.Object("referenceTime", referenceTime))

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = false AND
				i.status IN $acceptedStatuses AND
				i.techPaymentLinkRequestedAt IS NULL AND
				(i.updatedAt + duration({minutes: $delay}) < $referenceTime OR i.techInvoiceFinalizedSentAt + duration({minutes: $delay}) < $referenceTime)
			RETURN distinct(i), t.name limit 100`
	params := map[string]any{
		"delay":         minutesFromLastUpdate,
		"referenceTime": referenceTime,
		"now":           utils.Now(),
		"acceptedStatuses": []string{
			neo4jenum.InvoiceStatusDue.String(),
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
