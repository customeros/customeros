package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type InvoiceReadRepository interface {
	// GetInvoiceById returns the invoice node with the given id for tenant, error if not found or multiple found
	GetInvoiceById(ctx context.Context, tenant, invoiceId string) (*dbtype.Node, error)
	GetInvoices(ctx context.Context, tenant string, organizationId string, skip, limit int) (*utils.DbNodesWithTotalCount, error)
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

func (r *invoiceReadRepository) GetInvoices(ctx context.Context, tenant string, organizationId string, skip, limit int) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.GetInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("organizationId", organizationId))
	span.LogFields(log.Int("skip", skip))
	span.LogFields(log.Int("limit", limit))

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		countParams := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) "+
				" RETURN count(i) as count", tenant, tenant, tenant),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"skip":           skip,
			"limit":          limit,
		}

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			" MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) "+
				" RETURN i "+
				" order by i.createdAt desc "+
				" SKIP $skip LIMIT $limit", tenant, tenant, tenant),
			params)
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
	span.LogFields(log.String("invoiceId", invoiceId))

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
