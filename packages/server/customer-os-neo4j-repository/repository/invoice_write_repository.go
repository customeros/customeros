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
	ContractId      string                  `json:"contractId"`
	Currency        neo4jenum.Currency      `json:"currency"`
	DryRun          bool                    `json:"dryRun"`
	InvoiceNumber   string                  `json:"invoiceNumber"`
	PeriodStartDate time.Time               `json:"periodStartDate"`
	PeriodEndDate   time.Time               `json:"periodEndDate"`
	CreatedAt       time.Time               `json:"createdAt"`
	SourceFields    model.Source            `json:"sourceFields"`
	BillingCycle    neo4jenum.BillingCycle  `json:"billingCycle"`
	Status          neo4jenum.InvoiceStatus `json:"status"`
	Note            string                  `json:"note"`
}

type InvoiceFillFields struct {
	Amount                        float64                 `json:"amount"`
	VAT                           float64                 `json:"vat"`
	TotalAmount                   float64                 `json:"totalAmount"`
	UpdatedAt                     time.Time               `json:"updatedAt"`
	ContractId                    string                  `json:"contractId"`
	Currency                      neo4jenum.Currency      `json:"currency"`
	DryRun                        bool                    `json:"dryRun"`
	InvoiceNumber                 string                  `json:"invoiceNumber"`
	PeriodStartDate               time.Time               `json:"periodStartDate"`
	PeriodEndDate                 time.Time               `json:"periodEndDate"`
	BillingCycle                  neo4jenum.BillingCycle  `json:"billingCycle"`
	Status                        neo4jenum.InvoiceStatus `json:"status"`
	Note                          string                  `json:"note"`
	DomesticPaymentsBankInfo      string                  `json:"domesticPaymentsBankInfo"`
	InternationalPaymentsBankInfo string                  `json:"internationalPaymentsBankInfo"`
	CustomerName                  string                  `json:"customerName"`
	CustomerAddress               string                  `json:"customerAddress"`
	CustomerEmail                 string                  `json:"customerEmail"`
	ProviderLogoUrl               string                  `json:"providerLogoUrl"`
	ProviderName                  string                  `json:"providerName"`
	ProviderAddress               string                  `json:"providerAddress"`
}

type InvoiceUpdateFields struct {
	UpdatedAt    time.Time               `json:"updatedAt"`
	Status       neo4jenum.InvoiceStatus `json:"status"`
	UpdateStatus bool                    `json:"updateStatus"`
}

type InvoiceWriteRepository interface {
	CreateInvoiceForContract(ctx context.Context, tenant, invoiceId string, data InvoiceCreateFields) error
	FillInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceFillFields) error
	InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error
	SetInvoicePaymentRequested(ctx context.Context, tenant, invoiceId string) error
	UpdateInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceUpdateFields) error
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
								i.updatedAt=$updatedAt,
								i.status=$status
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
								i.periodEndDate=$periodEnd,
								i.billingCycle=$billingCycle,
								i.note=$note
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
		"billingCycle":  data.BillingCycle.String(),
		"status":        data.Status.String(),
		"note":          data.Note,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) FillInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceFillFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.FillInvoice")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})
							MERGE (t)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId}) 
							ON CREATE SET
								i:Invoice_%s,
								i.currency=$currency,
								i.dryRun=$dryRun,
								i.number=$number,
								i.currency=$currency,
								i.periodStartDate=$periodStart,
								i.periodEndDate=$periodEnd,
								i.billingCycle=$billingCycle,
								i.note=$note,
								i.domesticPaymentsBankInfo=$domesticPaymentsBankInfo,
								i.internationalPaymentsBankInfo=$internationalPaymentsBankInfo,
								i.customerName=$customerName,
								i.customerAddress=$customerAddress,
								i.customerEmail=$customerEmail,
								i.providerLogoUrl=$providerLogoUrl,
								i.providerName=$providerName,
								i.providerAddress=$providerAddress
							SET 
								i.updatedAt=$updatedAt,
								i.amount=$amount,
								i.vat=$vat,
								i.totalAmount=$totalAmount,
								i.status=$status,
								i.note=$note,
								i.domesticPaymentsBankInfo=$domesticPaymentsBankInfo,
								i.internationalPaymentsBankInfo=$internationalPaymentsBankInfo,
								i.customerName=$customerName,
								i.customerAddress=$customerAddress,
								i.customerEmail=$customerEmail,
								i.providerLogoUrl=$providerLogoUrl,
								i.providerName=$providerName,
								i.providerAddress=$providerAddress
							WITH c, i 
							MERGE (c)-[:HAS_INVOICE]->(i) 
							`, tenant)
	params := map[string]any{
		"tenant":                        tenant,
		"contractId":                    data.ContractId,
		"invoiceId":                     invoiceId,
		"updatedAt":                     data.UpdatedAt,
		"amount":                        data.Amount,
		"vat":                           data.VAT,
		"totalAmount":                   data.TotalAmount,
		"dryRun":                        data.DryRun,
		"number":                        data.InvoiceNumber,
		"currency":                      data.Currency.String(),
		"periodStart":                   data.PeriodStartDate,
		"periodEnd":                     data.PeriodEndDate,
		"billingCycle":                  data.BillingCycle.String(),
		"status":                        data.Status.String(),
		"note":                          data.Note,
		"domesticPaymentsBankInfo":      data.DomesticPaymentsBankInfo,
		"internationalPaymentsBankInfo": data.InternationalPaymentsBankInfo,
		"customerName":                  data.CustomerName,
		"customerAddress":               data.CustomerAddress,
		"customerEmail":                 data.CustomerEmail,
		"providerLogoUrl":               data.ProviderLogoUrl,
		"providerName":                  data.ProviderName,
		"providerAddress":               data.ProviderAddress,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) UpdateInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.UpdateInvoice")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
				SET i.updatedAt=$updatedAt`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
		"updatedAt": data.UpdatedAt,
	}
	if data.UpdateStatus && data.Status.String() != "" {
		cypher += `, i.status=$status`
		params["status"] = data.Status.String()
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
