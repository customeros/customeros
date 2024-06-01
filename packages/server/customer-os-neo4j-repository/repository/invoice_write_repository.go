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
	ContractId           string                  `json:"contractId"`
	Currency             neo4jenum.Currency      `json:"currency"`
	DryRun               bool                    `json:"dryRun"`
	OffCycle             bool                    `json:"offCycle"`
	Postpaid             bool                    `json:"postpaid"`
	Preview              bool                    `json:"preview"`
	PeriodStartDate      time.Time               `json:"periodStartDate"`
	PeriodEndDate        time.Time               `json:"periodEndDate"`
	CreatedAt            time.Time               `json:"createdAt"`
	IssuedDate           time.Time               `json:"issuedDate"`
	DueDate              time.Time               `json:"dueDate"`
	SourceFields         model.Source            `json:"sourceFields"`
	BillingCycleInMonths int64                   `json:"billingCycleInMonths"`
	Status               neo4jenum.InvoiceStatus `json:"status"`
	Note                 string                  `json:"note"`
}

type InvoiceFillFields struct {
	Amount                       float64                 `json:"amount"`
	VAT                          float64                 `json:"vat"`
	TotalAmount                  float64                 `json:"totalAmount"`
	UpdatedAt                    time.Time               `json:"updatedAt"`
	ContractId                   string                  `json:"contractId"`
	Currency                     neo4jenum.Currency      `json:"currency"`
	DryRun                       bool                    `json:"dryRun"`
	OffCycle                     bool                    `json:"offCycle"`
	Postpaid                     bool                    `json:"postpaid"`
	Preview                      bool                    `json:"preview"`
	InvoiceNumber                string                  `json:"invoiceNumber"`
	PeriodStartDate              time.Time               `json:"periodStartDate"`
	PeriodEndDate                time.Time               `json:"periodEndDate"`
	BillingCycleInMonths         int64                   `json:"billingCycleInMonths"`
	Status                       neo4jenum.InvoiceStatus `json:"status"`
	Note                         string                  `json:"note"`
	CustomerName                 string                  `json:"customerName"`
	CustomerEmail                string                  `json:"customerEmail"`
	CustomerAddressLine1         string                  `json:"customerAddressLine1"`
	CustomerAddressLine2         string                  `json:"customerAddressLine2"`
	CustomerAddressZip           string                  `json:"customerAddressZip"`
	CustomerAddressLocality      string                  `json:"customerAddressLocality"`
	CustomerAddressCountry       string                  `json:"customerAddressCountry"`
	CustomerAddressRegion        string                  `json:"customerAddressRegion"`
	ProviderLogoRepositoryFileId string                  `json:"providerLogoRepositoryFileId"`
	ProviderName                 string                  `json:"providerName"`
	ProviderEmail                string                  `json:"providerEmail"`
	ProviderAddressLine1         string                  `json:"providerAddressLine1"`
	ProviderAddressLine2         string                  `json:"providerAddressLine2"`
	ProviderAddressZip           string                  `json:"providerAddressZip"`
	ProviderAddressLocality      string                  `json:"providerAddressLocality"`
	ProviderAddressCountry       string                  `json:"providerAddressCountry"`
	ProviderAddressRegion        string                  `json:"providerAddressRegion"`
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
	UpdateInvoice(ctx context.Context, tenant, invoiceId string, data InvoiceUpdateFields) error
	MarkInvoiceFinalizedEventSent(ctx context.Context, tenant, invoiceId string) error
	MarkInvoiceFinalizedWebhookProcessed(ctx context.Context, tenant, invoiceId string) error
	MarkPayNotificationRequested(ctx context.Context, tenant, invoiceId string, requestedAt time.Time) error
	MarkPaymentLinkRequested(ctx context.Context, tenant, invoiceId string) error
	SetPaidInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error
	SetPayInvoiceNotificationSentAt(ctx context.Context, tenant, invoiceId string) error
	DeleteInitializedInvoice(ctx context.Context, tenant, invoiceId string) error
	DeletePreviewCycleInvoices(ctx context.Context, tenant, contractId, skipInvoiceId string) error
	DeletePreviewCycleInitializedInvoices(ctx context.Context, tenant, contractId, skipInvoiceId string) error
	DeleteDryRunInvoice(ctx context.Context, tenant, invoiceId string) error
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
								i.issuedDate=$issuedDate,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.dryRun=$dryRun,
								i.offCycle=$offCycle,
								i.postpaid=$postpaid,
								i.preview=$preview,
								i.currency=$currency,
								i.periodStartDate=$periodStart,
								i.periodEndDate=$periodEnd,
								i.billingCycleInMonths=$billingCycleInMonths,
								i.note=$note
							WITH c, i 
							MERGE (c)-[:HAS_INVOICE]->(i) 
							`, tenant)
	params := map[string]any{
		"tenant":               tenant,
		"contractId":           data.ContractId,
		"invoiceId":            invoiceId,
		"createdAt":            data.CreatedAt,
		"updatedAt":            data.CreatedAt,
		"dueDate":              utils.ToNeo4jDateAsAny(&data.DueDate),
		"issuedDate":           utils.ToDate(data.IssuedDate),
		"source":               data.SourceFields.Source,
		"sourceOfTruth":        data.SourceFields.Source,
		"appSource":            data.SourceFields.AppSource,
		"dryRun":               data.DryRun,
		"offCycle":             data.OffCycle,
		"postpaid":             data.Postpaid,
		"preview":              data.Preview,
		"currency":             data.Currency.String(),
		"periodStart":          utils.ToNeo4jDateAsAny(&data.PeriodStartDate),
		"periodEnd":            utils.ToNeo4jDateAsAny(&data.PeriodEndDate),
		"billingCycleInMonths": data.BillingCycleInMonths,
		"status":               data.Status.String(),
		"note":                 data.Note,
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
								i.offCycle=$offCycle,
								i.postpaid=$postpaid,
								i.preview=$preview,
								i.periodStartDate=$periodStart,
								i.periodEndDate=$periodEnd,
								i.billingCycleInMonths=$billingCycleInMonths
							SET 
								i.updatedAt=$updatedAt,
								i.number=$number,
								i.amount=$amount,
								i.vat=$vat,
								i.totalAmount=$totalAmount,
								i.status=$status,
								i.note=$note,
								i.customerName=$customerName,
								i.customerEmail=$customerEmail,
								i.customerAddressLine1=$customerAddressLine1,
								i.customerAddressLine2=$customerAddressLine2,
								i.customerAddressZip=$customerAddressZip,
								i.customerAddressLocality=$customerAddressLocality,
								i.customerAddressCountry=$customerAddressCountry,
								i.customerAddressRegion=$customerAddressRegion,
								i.providerLogoRepositoryFileId=$providerLogoRepositoryFileId,
								i.providerName=$providerName,
								i.providerEmail=$providerEmail,
								i.providerAddressLine1=$providerAddressLine1,
								i.providerAddressLine2=$providerAddressLine2,
								i.providerAddressZip=$providerAddressZip,
								i.providerAddressLocality=$providerAddressLocality,
								i.providerAddressCountry=$providerAddressCountry,
								i.providerAddressRegion=$providerAddressRegion
							WITH c, i 
							MERGE (c)-[:HAS_INVOICE]->(i) 
							`, tenant)
	params := map[string]any{
		"tenant":                       tenant,
		"contractId":                   data.ContractId,
		"invoiceId":                    invoiceId,
		"updatedAt":                    data.UpdatedAt,
		"amount":                       data.Amount,
		"vat":                          data.VAT,
		"totalAmount":                  data.TotalAmount,
		"dryRun":                       data.DryRun,
		"offCycle":                     data.OffCycle,
		"postpaid":                     data.Postpaid,
		"preview":                      data.Preview,
		"number":                       data.InvoiceNumber,
		"currency":                     data.Currency.String(),
		"periodStart":                  utils.ToNeo4jDateAsAny(&data.PeriodStartDate),
		"periodEnd":                    utils.ToNeo4jDateAsAny(&data.PeriodEndDate),
		"billingCycleInMonths":         data.BillingCycleInMonths,
		"status":                       data.Status.String(),
		"note":                         data.Note,
		"customerName":                 data.CustomerName,
		"customerEmail":                data.CustomerEmail,
		"customerAddressLine1":         data.CustomerAddressLine1,
		"customerAddressLine2":         data.CustomerAddressLine2,
		"customerAddressZip":           data.CustomerAddressZip,
		"customerAddressLocality":      data.CustomerAddressLocality,
		"customerAddressCountry":       data.CustomerAddressCountry,
		"customerAddressRegion":        data.CustomerAddressRegion,
		"providerLogoRepositoryFileId": data.ProviderLogoRepositoryFileId,
		"providerName":                 data.ProviderName,
		"providerEmail":                data.ProviderEmail,
		"providerAddressLine1":         data.ProviderAddressLine1,
		"providerAddressLine2":         data.ProviderAddressLine2,
		"providerAddressZip":           data.ProviderAddressZip,
		"providerAddressLocality":      data.ProviderAddressLocality,
		"providerAddressCountry":       data.ProviderAddressCountry,
		"providerAddressRegion":        data.ProviderAddressRegion,
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

func (r *invoiceWriteRepository) MarkInvoiceFinalizedEventSent(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.MarkInvoiceFinalizedEventSent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techInvoiceFinalizedSentAt=$now`, tenant)
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

func (r *invoiceWriteRepository) MarkInvoiceFinalizedWebhookProcessed(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.MarkInvoiceFinalizedWebhookProcessed")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techInvoiceFinalizedWebhookProcessedAt=$now`, tenant)
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

func (r *invoiceWriteRepository) DeleteInitializedInvoice(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.DeleteInitializedInvoice")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[r1:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})<-[r2:HAS_INVOICE]-(:Contract)
							WHERE i:Invoice_%s AND i.status=$initializedStatus
							DELETE r1,r2,i`, tenant)
	params := map[string]any{
		"tenant":            tenant,
		"invoiceId":         invoiceId,
		"initializedStatus": neo4jenum.InvoiceStatusInitialized.String(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *invoiceWriteRepository) DeletePreviewCycleInvoices(ctx context.Context, tenant, contractId, skipInvoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.DeletePreviewCycleInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId), log.String("skipInvoiceId", skipInvoiceId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract{id:$contractId})-[:HAS_INVOICE]->(i:Invoice {dryRun:true, preview: true, offCycle: false})
							   WHERE i.id <> $skipInvoiceId
			   OPTIONAL MATCH (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine) 
			   DETACH DELETE i, il`
	params := map[string]any{
		"tenant":        tenant,
		"contractId":    contractId,
		"skipInvoiceId": skipInvoiceId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *invoiceWriteRepository) DeletePreviewCycleInitializedInvoices(ctx context.Context, tenant, contractId, skipInvoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.DeletePreviewCycleInitializedInvoices")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("contractId", contractId), log.String("skipInvoiceId", skipInvoiceId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract{id:$contractId})-[:HAS_INVOICE]->(i:Invoice {dryRun:true, preview: true, offCycle: false})
							   WHERE i.id <> $skipInvoiceId and i.status=$initializedStatus
			   OPTIONAL MATCH (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine) 
			   DETACH DELETE i, il`
	params := map[string]any{
		"tenant":            tenant,
		"contractId":        contractId,
		"skipInvoiceId":     skipInvoiceId,
		"initializedStatus": neo4jenum.InvoiceStatusInitialized.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *invoiceWriteRepository) MarkPaymentLinkRequested(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceWriteRepository.MarkPaymentLinkRequested")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId})
							WHERE i:Invoice_%s
							SET i.techPaymentLinkRequestedAt=$now`, tenant)
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

func (r *invoiceWriteRepository) DeleteDryRunInvoice(ctx context.Context, tenant, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceReadRepository.DeleteDryRunInvoice")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("invoiceId", invoiceId))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$invoiceId, dryRun:true})
			   OPTIONAL MATCH (i)-[:HAS_INVOICE_LINE]->(il:InvoiceLine) 
			   DETACH DELETE i, il`
	params := map[string]any{
		"tenant":    tenant,
		"invoiceId": invoiceId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
