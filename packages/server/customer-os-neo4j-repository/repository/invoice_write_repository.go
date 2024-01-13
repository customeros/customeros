package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InvoiceWriteRepository interface {
	InvoiceNew(ctx context.Context, tenant, organizationId, id string, date, dueDate time.Time, dryRun bool, source, appSource string, createdAt time.Time) error
	InvoiceFill(ctx context.Context, tenant, id string, amount, vat, total float64, updatedAt time.Time) error
	InvoiceFillInvoiceLine(ctx context.Context, tenant, id string, index int64, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error
	InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error
}

type invoiceWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoiceWriteRepository(driver *neo4j.DriverWithContext, database string) InvoiceWriteRepository {
	return &invoiceWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoiceWriteRepository) InvoiceNew(ctx context.Context, tenant, organizationId, id string, date, dueDate time.Time, dryRun bool, source, appSource string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.InvoiceNew")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s {id:$organizationId})
							MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) 
							ON CREATE SET 
								i:Invoice_%s,
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.date=$date,
								i.dueDate=$dueDate,
								i.dryRun=$dryRun,
								i.amount=0.0,
								i.vat=0.0,
								i.total=0.0,
								i.pdfGenerated=false
							WITH o, i 
							MERGE (o)-[:HAS_INVOICE]->(i) 
							`, tenant, tenant)
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"id":             id,
		"createdAt":      createdAt,
		"updatedAt":      createdAt,
		"source":         source,
		"sourceOfTruth":  source,
		"appSource":      appSource,
		"dryRun":         dryRun,
		"date":           date,
		"dueDate":        dueDate,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) InvoiceFill(ctx context.Context, tenant, id string, amount, vat, total float64, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.InvoiceFill")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) 
							SET i.updatedAt=$updatedAt,
								i.amount=$amount,
								i.vat=$vat,
								i.total=$total
	`
	params := map[string]any{
		"tenant":    tenant,
		"id":        id,
		"updatedAt": updatedAt,
		"amount":    amount,
		"vat":       vat,
		"total":     total,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) InvoiceFillInvoiceLine(ctx context.Context, tenant, id string, index int64, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice_%s {id:$id})
							MERGE (i)<-[:INVOICE_LINE_BELONGS_TO_INVOICE]-(il:InvoiceLine {id:randomUUID()})
							ON CREATE SET 
								il:InvoiceLine_%s,
								il.createdAt=$createdAt,
								il.updatedAt=$updatedAt,
								il.index=$index,
								il.name=$name,
								il.price=$price,
								il.quantity=$quantity,
								il.amount=$amount,
								il.vat=$vat,
								il.total=$total
`, tenant, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"id":        id,
		"createdAt": createdAt,
		"updatedAt": createdAt,
		"index":     index,
		"name":      name,
		"price":     price,
		"quantity":  quantity,
		"amount":    amount,
		"vat":       vat,
		"total":     total,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice_%s {id:$id}) 
							SET 
								i.pdfGenerated=true, 
								i.repositoryFileId=$repositoryFileId, 
								i.updatedAt=$updatedAt`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"id":               id,
		"updatedAt":        updatedAt,
		"repositoryFileId": repositoryFileId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
