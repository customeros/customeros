package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InvoiceCreateFields struct {
	ContractId      string             `json:"contractId"`
	Currency        neo4jenum.Currency `json:"currency"`
	DryRun          bool               `json:"dryRun"`
	InvoiceNumber   string             `json:"invoiceNumber"`
	PeriodStartDate time.Time          `json:"periodStartDate"`
	PeriodEndDate   time.Time          `json:"periodEndDate"`
	CreatedAt       time.Time          `json:"createdAt"`
	SourceFields    model.Source       `json:"sourceFields"`
}

type InvoiceWriteRepository interface {
	CreateInvoiceForContract(ctx context.Context, tenant, invoiceId string, data InvoiceCreateFields) error
	InvoiceFill(ctx context.Context, tenant, id string, amount, vat, total float64, updatedAt time.Time) error
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

func (r *invoiceWriteRepository) CreateInvoiceForContract(ctx context.Context, tenant, invoiceId string, data InvoiceCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.CreateInvoiceForContract")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId}) 
							ON CREATE SET
								i.updatedAt=$updatedAt
							SET 
								i:Invoice_%s,
								i.createdAt=$createdAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.dryRun=$dryRun,
								i.number=$number,
								i.currency=$currency,
								i.periodStartDate=$periodStart,
								i.periodEndDate=$periodEnd
							WITH c, i 
							MERGE (c)-[:HAS_INVOICE]->(i) 
							`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    data.ContractId,
		"invoiceId":     invoiceId,
		"createdAt":     data.CreatedAt,
		"updatedAt":     data.CreatedAt,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.Source,
		"appSource":     data.SourceFields.AppSource,
		"dryRun":        data.DryRun,
		"number":        data.InvoiceNumber,
		"currency":      data.Currency.String(),
		"periodStart":   data.PeriodStartDate,
		"periodEnd":     data.PeriodEndDate,
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

func (r *invoiceWriteRepository) InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice_%s {id:$id}) 
							SET 
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
