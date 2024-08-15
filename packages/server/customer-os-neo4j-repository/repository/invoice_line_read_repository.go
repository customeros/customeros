package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type InvoiceLineReadRepository interface {
	GetAllForInvoice(ctx context.Context, tenant string, invoiceId string) ([]*dbtype.Node, error)
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

func (r *invoiceLineReadRepository) GetAllForInvoice(ctx context.Context, tenant string, invoiceId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceLineReadRepository.GetAllForInvoice")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("invoiceId", invoiceId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice{id:$invoiceId})-[:HAS_INVOICE_LINE]->(il:InvoiceLine)
		 RETURN il ORDER BY il.name`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
	}
	span.LogFields(log.String("cypher", cypher))
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
	return result.([]*dbtype.Node), err
}

func (r *invoiceLineReadRepository) GetAllForInvoices(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceLineReadRepository.GetAllForInvoice")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("ids", ids))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)-[:HAS_INVOICE_LINE]->(il:InvoiceLine)
		 WHERE i.id IN $ids 
		 RETURN il, i.id ORDER BY il.createdAt asc`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	span.LogFields(log.String("cypher", cypher))
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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *invoiceLineReadRepository) GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx context.Context, tenant, sliParentId string) (*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("sliParentId", sliParentId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)-[:HAS_INVOICE_LINE]->(il:InvoiceLine)-[:INVOICED]->(sli:ServiceLineItem {parentId:$parentId})
		WHERE NOT i.status IN $skipStatuses AND i.dryRun = false
		 RETURN il, i.id ORDER BY il.createdAt desc limit 1`
	params := map[string]any{
		"tenant":       tenant,
		"parentId":     sliParentId,
		"skipStatuses": []string{neo4jenum.InvoiceStatusInitialized.String(), neo4jenum.InvoiceStatusVoid.String()},
	}
	span.LogFields(log.String("cypher", cypher))
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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId)[0], err
}
