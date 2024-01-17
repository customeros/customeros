package repository

import (
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
