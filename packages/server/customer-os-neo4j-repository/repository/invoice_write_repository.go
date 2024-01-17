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
	InvoiceNew(ctx context.Context, tenant, contractId, id string, dryRun bool, number string, date, dueDate time.Time, source, appSource string, createdAt time.Time) error
	InvoiceFill(ctx context.Context, tenant, id string, amount, vat, total float64, updatedAt time.Time) error
	InvoiceFillInvoiceLine(ctx context.Context, tenant, id string, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error
	InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error
	SetInvoicePaymentRequested(ctx context.Context, tenant, invoiceId string) error
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

func (r *invoiceWriteRepository) InvoiceNew(ctx context.Context, tenant, contractId, id string, dryRun bool, number string, date, dueDate time.Time, source, appSource string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.InvoiceNew")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract_%s {id:$contractId})
							MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) 
							ON CREATE SET 
								i:Invoice_%s,
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.dryRun=$dryRun,
								i.number=$number,
								i.date=$date,
								i.dueDate=$dueDate,
								i.amount=0.0,
								i.vat=0.0,
								i.total=0.0,
								i.pdfGenerated=false
							WITH c, i 
							MERGE (c)-[:HAS_INVOICE]->(i) 
							`, tenant, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"id":            id,
		"createdAt":     createdAt,
		"updatedAt":     createdAt,
		"source":        source,
		"sourceOfTruth": source,
		"appSource":     appSource,
		"dryRun":        dryRun,
		"number":        number,
		"date":          date,
		"dueDate":       dueDate,
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

func (r *invoiceWriteRepository) InvoiceFillInvoiceLine(ctx context.Context, tenant, id string, name string, price float64, quantity int64, amount, vat, total float64, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id})
							WHERE i:Invoice_%s
							MERGE (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine {id:randomUUID()})
							ON CREATE SET 
								il:InvoiceLine_%s,
								il.createdAt=$createdAt,
								il.updatedAt=$updatedAt,
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

func (r *invoiceWriteRepository) SetInvoicePaymentRequested(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.SetInvoicePaymentRequested")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techPaymentRequestedAt=$now`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
		"now":       utils.Now(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
