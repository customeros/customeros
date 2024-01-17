package invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

type RequestBodyInvoiceReady struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	StripeCustomerId             string `json:"stripeCustomerId"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
}

type InvoiceEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	cfg          *config.EventNotifications
}

func NewInvoiceEventHandler(log logger.Logger, repositories *repository.Repositories, cfg *config.EventNotifications) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:          log,
		repositories: repositories,
		cfg:          cfg,
	}
}

func (h *InvoiceEventHandler) onInvoiceNewV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceNewV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceNewEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	// get currency
	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting contract %s: %s", eventData.ContractId, err.Error())
		return err
	}
	if contractDbNode == nil {
		err = errors.Errorf("Contract %s not found", eventData.ContractId)
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting contract %s: %s", eventData.ContractId, err.Error())
		return err
	}
	//contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	// TODO use currency from above contract. For now, use default currency
	// TODO temp code starts here to defaut tenant currency
	tenantSettingsDbNode, err := h.repositories.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting tenant settings for tenant %s: %s", eventData.Tenant, err.Error())
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)
	currency := tenantSettingsEntity.DefaultCurrency.String()
	if currency == "" {
		currency = neo4jenum.CurrencyUSD.String()
	}
	// TODO temp code ends here

	// fire fill invoice event

	return nil
}

func (h *InvoiceEventHandler) onInvoicePdfGeneratedV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePdfGeneratedV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePdfGeneratedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	invoiceDbNode, err := h.repositories.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting invoice %s: %s", invoiceId, err.Error())
		return err
	}
	if invoiceDbNode == nil {
		err = errors.Errorf("Invoice %s not found", invoiceId)
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting invoice %s: %s", invoiceId, err.Error())
		return err
	}
	invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode)

	err = h.invokeInvoiceReadyWebhook(ctx, eventData.Tenant, *invoiceEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error invoking invoice ready webhook for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) invokeInvoiceReadyWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.invokeInvoiceReadyWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	if invoice.DryRun {
		return nil
	}
	if h.cfg.EndPoints.InvoiceReady == "" {
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get stripe customer id for organization
	stripeCustomerId, err := h.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetFirstExternalIdForLinkedEntity(ctx, tenant, neo4jenum.Stripe.String(), organizationEntity.ID, neo4jutil.NodeLabelOrganization)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error getting stripe customer id for organization %s: %s", organizationEntity.ID, err.Error())
		return err
	}

	// convert amount to the smallest currency unit
	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.Amount)
	if err != nil {
		return fmt.Errorf("error converting amount to smallest currency unit: %v", err.Error())
	}

	requestBody := RequestBodyInvoiceReady{
		Tenant:                       tenant,
		Currency:                     invoice.Currency.String(),
		AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
		StripeCustomerId:             stripeCustomerId,
		InvoiceId:                    invoice.Id,
		InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", h.cfg.EndPoints.InvoiceReady, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %s", resp.Status)
	}

	// Request was successful
	err = h.repositories.Neo4jRepositories.InvoiceWriteRepository.SetInvoicePaymentRequested(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error setting invoice payment requested for invoice %s: %s", invoice.Id, err.Error())
	}

	return nil
}
