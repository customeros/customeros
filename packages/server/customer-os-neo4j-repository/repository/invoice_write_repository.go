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
	OffCycle        bool                    `json:"offCycle"`
	PeriodStartDate time.Time               `json:"periodStartDate"`
	PeriodEndDate   time.Time               `json:"periodEndDate"`
	CreatedAt       time.Time               `json:"createdAt"`
	DueDate         time.Time               `json:"dueDate"`
	SourceFields    model.Source            `json:"sourceFields"`
	BillingCycle    neo4jenum.BillingCycle  `json:"billingCycle"`
	Status          neo4jenum.InvoiceStatus `json:"status"`
	Note            string                  `json:"note"`
}

type InvoiceFillFields struct {
	Amount                        float64                 `json:"amount"`
	VAT                           float64                 `json:"vat"`
	SubtotalAmount                float64                 `json:"subtotalAmount"`
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
	CustomerEmail                 string                  `json:"customerEmail"`
	CustomerAddressLine1          string                  `json:"customerAddressLine1"`
	CustomerAddressLine2          string                  `json:"customerAddressLine2"`
	CustomerAddressZip            string                  `json:"customerAddressZip"`
	CustomerAddressLocality       string                  `json:"customerAddressLocality"`
	CustomerAddressCountry        string                  `json:"customerAddressCountry"`
	ProviderLogoUrl               string                  `json:"providerLogoUrl"`
	ProviderName                  string                  `json:"providerName"`
	ProviderEmail                 string                  `json:"providerEmail"`
	ProviderAddressLine1          string                  `json:"providerAddressLine1"`
	ProviderAddressLine2          string                  `json:"providerAddressLine2"`
	ProviderAddressZip            string                  `json:"providerAddressZip"`
	ProviderAddressLocality       string                  `json:"providerAddressLocality"`
	ProviderAddressCountry        string                  `json:"providerAddressCountry"`
}

type InvoiceUpdateFields struct {
	UpdatedAt         time.Time               `json:"updatedAt"`
	Status            neo4jenum.InvoiceStatus `json:"status"`
	PaymentLink       string                  `json:"paymentLink"`
	UpdateStatus      bool                    `json:"updateStatus"`
	UpdatePaymentLink bool                    `json:"updatePaymentLink"`
}

type InvoiceWriteRepository interface {
	CreateInvoiceForContract(ctx context.Context, tenant, invoiceId string, data InvoiceCreateFields) error
	FillInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceFillFields) error
	InvoicePdfGenerated(ctx context.Context, tenant, id, repositoryFileId string, updatedAt time.Time) error
	SetInvoicePaymentRequested(ctx context.Context, tenant, invoiceId string) error
	UpdateInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceUpdateFields) error
	MarkPayNotificationRequested(ctx context.Context, tenant, invoiceId string, requestedAt time.Time) error
	SetPaidInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error
	SetPayInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error
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
								i.dueDate=$dueDate,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.dryRun=$dryRun,
								i.offCycle=$offCycle,
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
		"dueDate":       data.DueDate,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.Source,
		"appSource":     data.SourceFields.AppSource,
		"dryRun":        data.DryRun,
		"offCycle":      data.OffCycle,
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
								i.currency=$currency,
								i.periodStartDate=$periodStart,
								i.periodEndDate=$periodEnd,
								i.billingCycle=$billingCycle,
								i.amount=$amount,
								i.vat=$vat,
								i.subtotalAmount=$subtotalAmount,
								i.totalAmount=$totalAmount,
								i.status=$status,
								i.note=$note,
								i.domesticPaymentsBankInfo=$domesticPaymentsBankInfo,
								i.internationalPaymentsBankInfo=$internationalPaymentsBankInfo,
								i.customerName=$customerName,
								i.customerEmail=$customerEmail,
								i.customerAddressLine1=$customerAddressLine1,
								i.customerAddressLine2=$customerAddressLine2,
								i.customerAddressZip=$customerAddressZip,
								i.customerAddressLocality=$customerAddressLocality,
								i.customerAddressCountry=$customerAddressCountry,
								i.providerLogoUrl=$providerLogoUrl,
								i.providerName=$providerName,
								i.providerEmail=$providerEmail,
								i.providerAddressLine1=$providerAddressLine1,
								i.providerAddressLine2=$providerAddressLine2,
								i.providerAddressZip=$providerAddressZip,
								i.providerAddressLocality=$providerAddressLocality,
								i.providerAddressCountry=$providerAddressCountry
							SET 
								i.updatedAt=$updatedAt,
								i.number=$number,
								i.amount=$amount,
								i.vat=$vat,
								i.subtotalAmount=$subtotalAmount,
								i.totalAmount=$totalAmount,
								i.status=$status,
								i.note=$note,
								i.domesticPaymentsBankInfo=$domesticPaymentsBankInfo,
								i.internationalPaymentsBankInfo=$internationalPaymentsBankInfo,
								i.customerName=$customerName,
								i.customerEmail=$customerEmail,
								i.customerAddressLine1=$customerAddressLine1,
								i.customerAddressLine2=$customerAddressLine2,
								i.customerAddressZip=$customerAddressZip,
								i.customerAddressLocality=$customerAddressLocality,
								i.customerAddressCountry=$customerAddressCountry,
								i.providerLogoUrl=$providerLogoUrl,
								i.providerName=$providerName,
								i.providerEmail=$providerEmail,
								i.providerAddressLine1=$providerAddressLine1,
								i.providerAddressLine2=$providerAddressLine2,
								i.providerAddressZip=$providerAddressZip,
								i.providerAddressLocality=$providerAddressLocality,
								i.providerAddressCountry=$providerAddressCountry
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
		"subtotalAmount":                data.SubtotalAmount,
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
		"customerEmail":                 data.CustomerEmail,
		"customerAddressLine1":          data.CustomerAddressLine1,
		"customerAddressLine2":          data.CustomerAddressLine2,
		"customerAddressZip":            data.CustomerAddressZip,
		"customerAddressLocality":       data.CustomerAddressLocality,
		"customerAddressCountry":        data.CustomerAddressCountry,
		"providerLogoUrl":               data.ProviderLogoUrl,
		"providerName":                  data.ProviderName,
		"providerEmail":                 data.ProviderEmail,
		"providerAddressLine1":          data.ProviderAddressLine1,
		"providerAddressLine2":          data.ProviderAddressLine2,
		"providerAddressZip":            data.ProviderAddressZip,
		"providerAddressLocality":       data.ProviderAddressLocality,
		"providerAddressCountry":        data.ProviderAddressCountry,
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
	if data.UpdatePaymentLink {
		cypher += `, i.paymentLink=$paymentLink`
		params["paymentLink"] = data.PaymentLink
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

func (r *invoiceWriteRepository) MarkPayNotificationRequested(ctx context.Context, tenant, invoiceId string, requestedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.MarkPayNotificationRequested")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
				SET i.techPayNotificationRequestedAt=$requestedAt`
	params := map[string]any{
		"tenant":      tenant,
		"invoiceId":   invoiceId,
		"requestedAt": requestedAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) SetPaidInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techPaidInvoiceNotificationSentAt=$now`, tenant)
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

func (r *invoiceWriteRepository) SetPayInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.SetPayInvoiceNotificationSentAt")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techPayInvoiceNotificationSentAt=$now`, tenant)
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
