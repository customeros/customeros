package neo4j_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type InvoiceLineReadRepository interface {
	GetAllForInvoice(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, invoiceId string) ([]*dbtype.Node, error)
	GetAllForInvoices(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx context.Context, tenant, sliParentId string) (*utils.DbNodeAndId, error)
}

type invoiceLineReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoiceLineReadRepository(driver *neo4j.DriverWithContext, database string) InvoiceLineReadRepository {
	return &invoiceLineReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoiceLineReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *invoiceLineReadRepository) GetAllForInvoice(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, invoiceId string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceLineReadRepository.GetAllForInvoice")
	defer spans.Finish()

	spans.LogKV("invoiceId", invoiceId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice{id:$invoiceId})-[:HAS_INVOICE_LINE]->(il:InvoiceLine)
		 RETURN il ORDER BY il.name`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *invoiceLineReadRepository) GetAllForInvoices(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceLineReadRepository.GetAllForInvoice")
	defer spans.Finish()

	spans.LogObjectAsJson("ids", ids)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)-[:HAS_INVOICE_LINE]->(il:InvoiceLine)
		 WHERE i.id IN $ids 
		 RETURN il, i.id ORDER BY il.createdAt asc`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *invoiceLineReadRepository) GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx context.Context, tenant, sliParentId string) (*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId")
	defer spans.Finish()

	spans.LogKV("sliParentId", sliParentId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)-[:HAS_INVOICE_LINE]->(il:InvoiceLine)-[:INVOICED]->(sli:ServiceLineItem {parentId:$parentId})
		WHERE NOT i.status IN $skipStatuses AND i.dryRun = false
		 RETURN il, i.id ORDER BY il.createdAt desc limit 1`
	params := map[string]any{
		"tenant":       tenant,
		"parentId":     sliParentId,
		"skipStatuses": []string{neo4jenum.InvoiceStatusInitialized.String()},
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId)[0], err
}
