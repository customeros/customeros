package invoice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/webhook"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type eventMetadata struct {
	UserId string `json:"user-id"`
}

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

type InvoiceEventHandler struct {
	log         logger.Logger
	cfg         config.Config
	grpcClients *grpc_client.Clients
	services    *service.Services
}

func NewInvoiceEventHandler(log logger.Logger, services *service.Services, cfg config.Config, grpcClients *grpc_client.Clients) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:         log,
		services:    services,
		cfg:         cfg,
		grpcClients: grpcClients,
	}
}

func (h *InvoiceEventHandler) onInvoiceFillRequestedV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceFillRequestedV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceFillRequestedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	invoiceEntity, err := h.services.CommonServices.InvoiceService.GetById(ctx, invoiceId)
	if err != nil {
		return err
	}

	if invoiceEntity.OffCycle {

		sliDbNodes, err := h.services.CommonServices.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForContract(ctx, eventData.Tenant, eventData.ContractId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ServiceLineItemReadRepository.GetServiceLineItemsForContract"))
			h.log.Errorf("Error getting service line items for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}

		var sliEntities neo4jentity.ServiceLineItemEntities
		for _, sliDbNode := range sliDbNodes {
			sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
			if sliEntity != nil {
				sliEntities = append(sliEntities, *neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode))
			}
		}

		invoiceEntity, invoiceLines, err := h.services.CommonServices.InvoiceService.FillOffCyclePrepaidInvoice(ctx, invoiceEntity, sliEntities)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "InvoiceService.FillOffCyclePrepaidInvoice"))
			h.log.Errorf("Error filling invoice %s: %s", invoiceId, err.Error())
			return err
		}

		allInvoiceLinesEmpty := true
		for _, invoiceLine := range invoiceLines {
			if invoiceLine.Amount != 0 {
				allInvoiceLinesEmpty = false
				break
			}
		}

		if allInvoiceLinesEmpty {
			_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
				return h.grpcClients.InvoiceClient.PermanentlyDeleteInitializedInvoice(ctx, &invoicepb.PermanentlyDeleteInitializedInvoiceRequest{
					Tenant:    eventData.Tenant,
					InvoiceId: invoiceEntity.Id,
					AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "PermanentlyDeleteInitializedInvoice"))
				h.log.Errorf("Error permanently deleting draft invoice {%s}: {%s}", invoiceEntity.Id, err.Error())
			}
			return err
		} else {
			return h.prepareAndCallFillInvoice(ctx, eventData.Tenant, eventData.ContractId, *invoiceEntity, invoiceEntity.Amount, invoiceEntity.Vat, invoiceEntity.TotalAmount, invoiceLines, span)
		}
	} else {

		sliDbNodes, err := h.services.CommonServices.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForContract(ctx, eventData.Tenant, eventData.ContractId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "ServiceLineItemReadRepository.GetServiceLineItemsForContract"))
			h.log.Errorf("Error getting service line items for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}

		var sliEntities neo4jentity.ServiceLineItemEntities
		for _, sliDbNode := range sliDbNodes {
			sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
			if sliEntity != nil {
				sliEntities = append(sliEntities, *neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode))
			}
		}

		invoiceEntity, invoiceLines, err := h.services.CommonServices.InvoiceService.FillCycleInvoice(ctx, invoiceEntity, sliEntities)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "InvoiceService.FillCycleInvoice"))
			h.log.Errorf("Error filling invoice %s: %s", invoiceId, err.Error())
			return err
		}

		return h.prepareAndCallFillInvoice(ctx, eventData.Tenant, eventData.ContractId, *invoiceEntity, invoiceEntity.Amount, invoiceEntity.Vat, invoiceEntity.TotalAmount, invoiceLines, span)
	}
}

func (h *InvoiceEventHandler) prepareAndCallFillInvoice(ctx context.Context, tenant string, contractId string, invoiceEntity neo4jentity.InvoiceEntity, amount, vat, totalAmount float64, invoiceLines []*invoicepb.InvoiceLine, span opentracing.Span) error {
	var contractEntity neo4jentity.ContractEntity
	var tenantSettingsEntity *neo4jentity.TenantSettingsEntity

	//load contract from neo4j
	contract, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetContractById"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetContractById")
	}
	if contract != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contract)
	} else {
		return errors.New("contract is nil")
	}

	//load tenant settings from neo4j
	tenantSettings, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	if tenantSettings != nil {
		tenantSettingsEntity = neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettings)
	} else {
		tracing.TraceErr(span, errors.New("tenantSettings is nil"))
		return errors.New("tenantSettings is nil")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, tenant, true)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return err
	}

	contractCountry := contractEntity.Country
	countryDbNode, _ := h.services.CommonServices.Neo4jRepositories.CountryReadRepository.GetCountryByCodeIfExists(ctx, contractCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		contractCountry = countryEntity.Name
	}
	tenantBillingProfileCountry := tenantBillingProfileEntity.Country
	countryDbNode, _ = h.services.CommonServices.Neo4jRepositories.CountryReadRepository.GetCountryByCodeIfExists(ctx, tenantBillingProfileCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		tenantBillingProfileCountry = countryEntity.Name
	}

	invoiceNumber := ""
	if !invoiceEntity.OffCycle {
		filledInvoiceDbNode, err := h.services.CommonServices.Neo4jRepositories.InvoiceReadRepository.GetFirstPreviewFilledInvoice(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "GetFirstPreviewFilledInvoice"))
		}
		if filledInvoiceDbNode != nil {
			filledInvoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(filledInvoiceDbNode)
			invoiceNumber = filledInvoiceEntity.Number
		}
	}

	err = h.callFillInvoice(ctx,
		tenant,
		invoiceEntity.Id,
		invoiceNumber,
		invoiceEntity.DryRun,
		invoiceEntity.Preview,
		contractEntity.ContractStatus,
		contractEntity.OrganizationLegalName,
		contractEntity.InvoiceEmail,
		contractEntity.AddressLine1, contractEntity.AddressLine2, contractEntity.Zip, contractEntity.Locality, contractCountry, contractEntity.Region,
		tenantSettingsEntity.LogoRepositoryFileId,
		tenantBillingProfileEntity.LegalName,
		tenantBillingProfileEntity.SendInvoicesFrom,
		tenantBillingProfileEntity.AddressLine1, tenantBillingProfileEntity.AddressLine2, tenantBillingProfileEntity.Zip, tenantBillingProfileEntity.Locality, tenantBillingProfileCountry, tenantBillingProfileEntity.Region,
		contractEntity.InvoiceNote,
		amount,
		vat,
		totalAmount,
		invoiceLines,
		span)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "callFillInvoice"))
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) callFillInvoice(ctx context.Context, tenant, invoiceId, invoiceNumber string, dryRun, preview bool, contractStatus neo4jenum.ContractStatus,
	customerName, customerEmail, customerAddressLine1, customerAddressLine2, customerAddressZip, customerAddressLocality, customerAddressCountry, customerAddressRegion,
	providerLogoRepositoryFileId, providerName, providerEmail, providerAddressLine1, providerAddressLine2, providerAddressZip, providerAddressLocality, providerAddressCountry, providerAddressRegion,
	note string, amount, vat, total float64, invoiceLines []*invoicepb.InvoiceLine, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	now := time.Now()

	invoiceStatus := invoicepb.InvoiceStatus_INVOICE_STATUS_DUE
	if len(invoiceLines) == 0 {
		invoiceStatus = invoicepb.InvoiceStatus_INVOICE_STATUS_EMPTY
	} else {
		if dryRun && preview {
			if contractStatus == neo4jenum.ContractStatusOutOfContract {
				invoiceStatus = invoicepb.InvoiceStatus_INVOICE_STATUS_ON_HOLD
			} else {
				invoiceStatus = invoicepb.InvoiceStatus_INVOICE_STATUS_SCHEDULED
			}
		} else if total == 0 {
			invoiceStatus = invoicepb.InvoiceStatus_INVOICE_STATUS_PAID
		}
	}
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return h.grpcClients.InvoiceClient.FillInvoice(ctx, &invoicepb.FillInvoiceRequest{
			Tenant:    tenant,
			InvoiceId: invoiceId,
			DryRun:    dryRun,
			Note:      note,
			Customer: &invoicepb.FillInvoiceCustomer{
				Name:         customerName,
				Email:        customerEmail,
				AddressLine1: customerAddressLine1,
				AddressLine2: customerAddressLine2,
				Zip:          customerAddressZip,
				Locality:     customerAddressLocality,
				Country:      customerAddressCountry,
				Region:       customerAddressRegion,
			},
			Provider: &invoicepb.FillInvoiceProvider{
				LogoRepositoryFileId: providerLogoRepositoryFileId,
				Name:                 providerName,
				Email:                providerEmail,
				AddressLine1:         providerAddressLine1,
				AddressLine2:         providerAddressLine2,
				Zip:                  providerAddressZip,
				Locality:             providerAddressLocality,
				Country:              providerAddressCountry,
				Region:               providerAddressRegion,
			},
			Amount:        amount,
			Vat:           vat,
			Total:         total,
			InvoiceLines:  invoiceLines,
			UpdatedAt:     utils.ConvertTimeToTimestampPtr(&now),
			AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
			Status:        invoiceStatus,
			InvoiceNumber: invoiceNumber,
		})
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "FillInvoice"))
		h.log.Errorf("Error sending the fill invoice request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) onInvoicePdfGeneratedV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePdfGeneratedV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePdfGeneratedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	invoiceEntity, err := h.services.CommonServices.InvoiceService.GetById(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceService.GetById"))
		return err
	}
	if invoiceEntity == nil {
		err = fmt.Errorf("invoice %s not found", invoiceId)
		tracing.TraceErr(span, errors.Wrap(err, "invoiceEntity == nil"))
		return err
	}

	if invoiceEntity.DryRun {
		return nil
	}

	err = h.slackInvoiceFinalizedWebhook(ctx, eventData.Tenant, *invoiceEntity)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "slackInvoiceFinalizedWebhook"))
		h.log.Errorf("error invoking slack invoice finalized webhook for invoice %s: %s", invoiceId, err.Error())
	}

	// do not dispatch invoice finalized event if it was already dispatched
	if invoiceEntity.InvoiceInternalFields.InvoiceFinalizedWebhookProcessedAt == nil {
		// dispatch invoice finalized event
		err = h.dispatchInvoiceFinalizedEvent(ctx, eventData.Tenant, *invoiceEntity)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "dispatchInvoiceFinalizedEvent"))
			h.log.Errorf("Error dispatching invoice finalized event for invoice %s: %s", invoiceId, err.Error())
			// TODO: must implement retry mechanism for dispatching invoice finalized event
		}
		err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, eventData.Tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyFinalizedWebhookProcessedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
			h.log.Errorf("Error setting invoice finalized webhook processed for invoice %s: %s", invoiceEntity.Id, err.Error())
		}
	}

	return nil
}

func (h *InvoiceEventHandler) slackInvoiceFinalizedWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.slackInvoiceFinalizedWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	if h.cfg.EventNotifications.SlackConfig.InternalAlertsRegisteredWebhook == "" {
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: fmt.Sprintf("Tenant %s, Invoice %s has been finalized for customer %s", tenant, invoice.Number, organizationEntity.Name)}
	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	// Send POST request
	resp, err := http.Post(h.cfg.EventNotifications.SlackConfig.InternalAlertsRegisteredWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	span.LogFields(log.String("result.status", resp.Status))

	return nil
}

func (h *InvoiceEventHandler) dispatchInvoiceFinalizedEvent(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.dispatchInvoiceFinalizedEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice to build payload for webhook
	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractsForOrganizations"))
		h.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	contractEntity := neo4jentity.ContractEntity{}
	if len(contractDbNode) > 0 && contractDbNode[0] != nil {
		node := contractDbNode[0].Node
		if node != nil {
			contractEntity = *neo4jmapper.MapDbNodeToContractEntity(node)
		}
	}

	// get invoice line items linked to invoice to build payload for webhook
	invoiceLineDbNodes, err := h.services.CommonServices.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoice(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetAllForInvoice"))
		h.log.Errorf("Error getting invoice line items for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	ilEntities := []*neo4jentity.InvoiceLineEntity{}
	for _, ilDbNode := range invoiceLineDbNodes {
		ilEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNode)
		ilEntities = append(ilEntities, ilEntity)
	}

	webhookPayload := webhook.PopulateInvoicePayload(&invoice, &organizationEntity, &contractEntity, ilEntities)
	// dispatch the event
	err = webhook.DispatchWebhook(
		ctx,
		tenant,
		webhook.WebhookEventInvoiceFinalized,
		webhookPayload,
		h.services.CommonServices.PostgresRepositories,
		h.cfg,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "DispatchWebhook"))
		h.log.Errorf("Error dispatching invoice finalized event for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) dispatchInvoicePaidEvent(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.dispatchInvoicePaidEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		h.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice to build payload for webhook
	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractsForOrganizations"))
		h.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	contractEntity := neo4jentity.ContractEntity{}
	if len(contractDbNode) > 0 && contractDbNode[0] != nil {
		node := contractDbNode[0].Node
		if node != nil {
			contractEntity = *neo4jmapper.MapDbNodeToContractEntity(node)
		}
	}

	// get invoice line items linked to invoice to build payload for webhook
	invoiceLineDbNodes, err := h.services.CommonServices.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoice(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetAllForInvoice"))
		h.log.Errorf("Error getting invoice line items for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	ilEntities := []*neo4jentity.InvoiceLineEntity{}
	for _, ilDbNode := range invoiceLineDbNodes {
		ilEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNode)
		ilEntities = append(ilEntities, ilEntity)
	}

	webhookPayload := webhook.PopulateInvoicePayload(&invoice, &organizationEntity, &contractEntity, ilEntities)
	// dispatch the event
	err = webhook.DispatchWebhook(
		ctx,
		tenant,
		webhook.WebhookEventInvoiceStatusPaid,
		webhookPayload,
		h.services.CommonServices.PostgresRepositories,
		h.cfg,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "DispatchWebhook"))
		h.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) generateInvoicePDFV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceSubscriber.generateInvoicePDFV1")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData invoice.InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	var contractEntity *neo4jentity.ContractEntity
	var invoiceEntity *neo4jentity.InvoiceEntity
	var invoiceLineEntities = []*neo4jentity.InvoiceLineEntity{}

	//load invoice
	invoiceEntity, err := h.services.CommonServices.InvoiceService.GetById(ctx, invoiceId)
	if err != nil {
		return err
	}

	// load contract
	contractNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	//load invoice lines
	invoiceLinesNodes, err := h.services.CommonServices.Neo4jRepositories.InvoiceLineReadRepository.GetAllForInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetAllForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetAllForInvoice")
	}
	if invoiceLinesNodes != nil {
		for _, invoiceLineNode := range invoiceLinesNodes {
			invoiceLineEntities = append(invoiceLineEntities, neo4jmapper.MapDbNodeToInvoiceLineEntity(invoiceLineNode))
		}
	} else {
		return errors.New("invoiceLinesNodes is nil")
	}

	invoiceHasVat := false

	if invoiceEntity.Vat > 0 {
		invoiceHasVat = true
	}

	dataForPdf := map[string]interface{}{
		"Tenant":                       eventData.Tenant,
		"CustomerName":                 invoiceEntity.Customer.Name,
		"CustomerEmail":                invoiceEntity.Customer.Email,
		"CustomerAddressLine1":         invoiceEntity.Customer.AddressLine1,
		"CustomerAddressLine2":         invoiceEntity.Customer.AddressLine2,
		"CustomerAddressLine3":         utils.JoinNonEmpty(", ", invoiceEntity.Customer.Locality, invoiceEntity.Customer.Zip),
		"CustomerAddressLine4":         utils.JoinNonEmpty(", ", invoiceEntity.Customer.Region, invoiceEntity.Customer.Country),
		"ProviderLogoExtension":        "",
		"ProviderLogoRepositoryFileId": invoiceEntity.Provider.LogoRepositoryFileId,
		"ProviderName":                 invoiceEntity.Provider.Name,
		"ProviderEmail":                invoiceEntity.Provider.Email,
		"ProviderAddressLine1":         invoiceEntity.Provider.AddressLine1,
		"ProviderAddressLine2":         invoiceEntity.Provider.AddressLine2,
		"ProviderAddressLine3":         utils.JoinNonEmpty(", ", invoiceEntity.Provider.Locality, invoiceEntity.Provider.Zip),
		"ProviderAddressLine4":         utils.JoinNonEmpty(", ", invoiceEntity.Provider.Region, invoiceEntity.Provider.Country),
		"InvoiceNumber":                invoiceEntity.Number,
		"InvoiceIssueDate":             invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		"InvoiceDueDate":               invoiceEntity.DueDate.Format("02 Jan 2006"),
		"InvoiceCurrency":              invoiceEntity.Currency.String() + "" + invoiceEntity.Currency.Symbol(),
		"InvoiceSubtotal":              utils.FormatAmount(invoiceEntity.Amount, 2),
		"InvoiceTotal":                 utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceAmountDue":             utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceLineItems":             []map[string]string{},
		"Note":                         invoiceEntity.Note,
		"CanPayByCheck":                contractEntity.Check,
		"DryRun":                       invoiceEntity.DryRun,
	}

	// Include bank details
	if contractEntity.CanPayWithBankTransfer {
		//load tenant billing profile from neo4j
		tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
			return err
		}
		if tenantBillingProfileEntity.CanPayWithBankTransfer {
			bankAccountDbNodes, err := h.services.CommonServices.Neo4jRepositories.BankAccountReadRepository.GetBankAccounts(ctx, eventData.Tenant)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetBankAccounts"))
				return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.GetBankAccounts")
			}
			for _, bankAccountDbNode := range bankAccountDbNodes {
				bankAccountEntity := neo4jmapper.MapDbNodeToBankAccountEntity(bankAccountDbNode)
				if bankAccountEntity.Currency == invoiceEntity.Currency {
					dataForPdf["BankDetailsAvailable"] = true
					dataForPdf["BankAccountName"] = bankAccountEntity.BankName
					dataForPdf["BankAccountNumber"] = bankAccountEntity.AccountNumber
					dataForPdf["BankAccountIBAN"] = bankAccountEntity.Iban
					dataForPdf["BankAccountBIC"] = bankAccountEntity.Bic
					dataForPdf["BankAccountSortCode"] = bankAccountEntity.SortCode
					dataForPdf["BankAccountRoutingNumber"] = bankAccountEntity.RoutingNumber
					dataForPdf["BankAccountOtherDetails"] = bankAccountEntity.OtherDetails
					break
				}
			}
		}
	}

	if invoiceHasVat {
		dataForPdf["InvoiceVat"] = fmt.Sprintf("%.2f", invoiceEntity.Vat)
	}

	for _, invoiceLine := range invoiceLineEntities {
		invoiceLineItem := map[string]string{
			"Name":      invoiceLine.Name,
			"Quantity":  fmt.Sprintf("%d", invoiceLine.Quantity),
			"UnitPrice": invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Price, 2),
			"Amount":    invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Amount, 2),
			"Vat":       invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Vat, 2),
		}
		sliDbNode, _ := h.services.CommonServices.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, invoiceLine.ServiceLineItemId)
		sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)

		if invoiceLine.BilledType == neo4jenum.BilledTypeOnce {
			invoiceLineItem["InvoiceLineSubtitle"] = sliEntity.StartedAt.Format("02 Jan 2006")
		}
		// if the invoice line item does not have a subtitle, we will use the period start and end date
		if _, ok := invoiceLineItem["InvoiceLineSubtitle"]; !ok {
			invoiceLineSubtitle := fmt.Sprintf("%s - %s", invoiceEntity.PeriodStartDate.Format("02 Jan 2006"), invoiceEntity.PeriodEndDate.Format("02 Jan 2006"))
			if sliEntity.Billed.IsRecurrent() && sliEntity.Billed.InMonths() != invoiceEntity.BillingCycleInMonths {
				invoiceLineSubtitle += ". "
				invoiceLineSubtitle += invoiceEntity.Currency.Symbol()
				invoiceLineSubtitle += utils.FormatAmount(sliEntity.Price, 2)
				invoiceLineSubtitle += "/"
				switch sliEntity.Billed {
				case neo4jenum.BilledTypeMonthly:
					invoiceLineSubtitle += "month"
				case neo4jenum.BilledTypeQuarterly:
					invoiceLineSubtitle += "quarter"
				case neo4jenum.BilledTypeAnnually:
					invoiceLineSubtitle += "year"
				}
			}
			invoiceLineItem["InvoiceLineSubtitle"] = invoiceLineSubtitle
		}

		if invoiceHasVat {
			invoiceLineItem["InvoiceHasVat"] = "true"
		}

		dataForPdf["InvoiceLineItems"] = append(dataForPdf["InvoiceLineItems"].([]map[string]string), invoiceLineItem)
	}

	//prepare the temp html file
	tmpInvoiceFile, err := os.CreateTemp("", "invoice_*.html")
	if err != nil {
		return errors.Wrap(err, "os.TempFile")
	}
	defer os.Remove(tmpInvoiceFile.Name()) // Delete the temporary HTML file when done
	defer tmpInvoiceFile.Close()

	if invoiceEntity.Provider.LogoRepositoryFileId != "" {
		fileMetadata, err := h.services.FileStoreApiService.GetFileMetadata(eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, span)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "GetFileMetadata"))
			h.log.Errorf("Error getting file metadata for file %s: %s", invoiceEntity.Provider.LogoRepositoryFileId, err.Error())
		} else {
			dataForPdf["ProviderLogoExtension"] = GetFileExtensionFromMetadata(fileMetadata)
		}
	}

	//fill the template with data and store it in temp
	err = FillInvoiceHtmlTemplate(ctx, tmpInvoiceFile, dataForPdf)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "FillInvoiceHtmlTemplate"))
		return errors.Wrap(err, "FillInvoiceHtmlTemplate")
	}

	//convert the temp to pdf
	// Max attempts
	maxAttempts := 3
	var pdfBytes *[]byte

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Try to convert the temp to pdf
		pdfBytes, err = ConvertInvoiceHtmlToPdf(ctx, h.services.FileStoreApiService, h.cfg.Subscriptions.InvoiceSubscription.PdfConverterUrl, tmpInvoiceFile, dataForPdf)
		if err == nil {
			// Success, no need to retry
			break
		}
		// Log the error and trace on failure
		tracing.TraceErr(span, errors.Wrap(err, "ConvertInvoiceHtmlToPdf"))
		if attempt == maxAttempts {
			// Return the error if it's the last attempt
			return errors.Wrap(err, "ConvertInvoiceHtmlToPdf")
		}
	}

	if pdfBytes == nil {
		return errors.New("pdfBytes is nil")
	}

	//TODO remove this at some point when we are sure that the pdf is generated correctly
	// Save the PDF file to disk
	os.WriteFile("output.pdf", *pdfBytes, 0644)

	basePath := fmt.Sprintf("/INVOICE/%d/%s", invoiceEntity.CreatedAt.Year(), invoiceEntity.CreatedAt.Format("01"))

	if invoiceEntity.DryRun {
		basePath = basePath + "/DRY_RUN"
	}

	fileDTO, err := h.services.FileStoreApiService.UploadSingleFileBytes(eventData.Tenant, basePath, invoiceEntity.Id, "Invoice - "+invoiceEntity.Number+".pdf", *pdfBytes, span)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UploadSingleFileBytes"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.UploadSingleFileBytes")
	}

	if fileDTO.Id == "" {
		return errors.New("fileDTO.Id is empty")
	}

	err = h.callPdfGeneratedInvoice(ctx, eventData.Tenant, invoiceId, fileDTO.Id, span)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "callPdfGeneratedInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.CallPdfGeneratedInvoice")
	}

	return nil
}

func (s *InvoiceEventHandler) callPdfGeneratedInvoice(ctx context.Context, tenant, invoiceId, repositoryFileId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.grpcClients.InvoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
			Tenant:           tenant,
			InvoiceId:        invoiceId,
			RepositoryFileId: repositoryFileId,
			AppSource:        constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "PdfGeneratedInvoice"))
		s.log.Errorf("Error sending the pdf generated request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) onInvoiceVoidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceVoidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceVoidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	invoiceEntity, err := h.services.CommonServices.InvoiceService.GetById(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceService.GetById"))
		return err
	}
	if invoiceEntity.DryRun {
		span.LogFields(log.String("result", "dry run, skipping"))
		return nil
	}

	invoiceLines, err := h.services.CommonServices.InvoiceService.GetInvoiceLinesForInvoices(ctx, []string{invoiceId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoiceLinesForInvoices"))
		return err
	}

	allInvoiceLinesEmpty := true
	for _, invoiceLine := range *invoiceLines {
		if invoiceLine.Amount != float64(0) {
			allInvoiceLinesEmpty = false
			break
		}
	}

	if invoiceEntity.DryRun || allInvoiceLinesEmpty {
		return nil
	}

	// void notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.VoidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load contract
	contractEntity := neo4jentity.ContractEntity{}
	contractNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return nil
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := commonService.PostmarkEmail{
		WorkflowId:    commonService.WorkflowInvoiceVoided,
		MessageStream: commonService.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            invoiceEntity.Customer.Email,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(commonService.WorkflowInvoiceVoidedSubject, invoiceEntity.Number), // "Voided invoice " + invoiceEntity.Number,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{issueDate}}":      invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		},
		Attachments: []commonService.PostmarkEmailAttachment{},
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.services.CommonServices.PostmarkService.SendNotification(ctx, postmarkEmail, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		h.log.Errorf("Error sending invoice voided notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.SetVoidInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetVoidInvoiceNotificationSentAt"))
		h.log.Errorf("Error setting invoice void notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoicePaidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePaidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePaidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	evtMetadata := eventMetadata{}
	if err := json.Unmarshal(evt.Metadata, &evtMetadata); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "json.Unmarshal"))
	}
	eventTriggeredByUser := evtMetadata.UserId != ""

	var invoiceEntity *neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	//load invoice
	invoiceEntity, err := h.services.CommonServices.InvoiceService.GetById(ctx, invoiceId)
	if err != nil {
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return nil
	}

	// do not dispatch invoice paid event if it was already dispatched
	if invoiceEntity.InvoiceInternalFields.InvoicePaidWebhookProcessedAt == nil {
		// dispatch invoice paid event
		err = h.dispatchInvoicePaidEvent(ctx, eventData.Tenant, *invoiceEntity)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "dispatchInvoicePaidEvent"))
			h.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoiceId, err.Error())
		} else {
			err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, eventData.Tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyPaidWebhookProcessedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
				h.log.Errorf("Error setting invoice paid webhook processed for invoice %s: %s", invoiceEntity.Id, err.Error())
			}
		}
	}

	// paid notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.PaidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load contract
	contractNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return nil
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return nil
	}
	tenantSettingsDbNode, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := commonService.PostmarkEmail{
		MessageStream: commonService.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentDate}}":    utils.Now().Format("02 Jan 2006"),
		},
		Attachments: []commonService.PostmarkEmailAttachment{},
	}
	if tenantSettingsEntity.StripeCustomerPortalLink != "" {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = fmt.Sprintf(`PS: If you pay by card you can manage your billing details <a href="%s">here</a>.`, tenantSettingsEntity.StripeCustomerPortalLink)
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = `PS: If you pay by card you can manage your billing details here.`
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = tenantSettingsEntity.StripeCustomerPortalLink
	} else {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = ""
	}

	if eventTriggeredByUser {
		postmarkEmail.WorkflowId = commonService.WorkflowInvoicePaymentReceived
		postmarkEmail.Subject = fmt.Sprintf(commonService.WorkflowInvoicePaymentReceivedSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	} else {
		postmarkEmail.WorkflowId = commonService.WorkflowInvoicePaid
		postmarkEmail.Subject = fmt.Sprintf(commonService.WorkflowInvoicePaidSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	}

	err = h.appendInvoiceFileToEmailAsAttachment(ctx, eventData.Tenant, *invoiceEntity, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendInvoiceFileToEmailAsAttachment"))
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.services.CommonServices.PostmarkService.SendNotification(ctx, postmarkEmail, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		h.log.Errorf("Error sending invoice paid notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPaidInvoiceNotificationSentAt"))
		h.log.Errorf("Error setting invoice paid notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoicePayNotificationV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoicePayNotificationV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePayNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice entity
	invoiceNode, err := h.services.CommonServices.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	// Do not send email if invoice is dry run or total amount is 0 or invoice is not due or overdue
	invoiceStatusAllowedForPayNotification := invoiceEntity.IsDue() || invoiceEntity.IsOverdue()
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceStatusAllowedForPayNotification {
		span.LogFields(log.String("result", "skipped pay notification"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		h.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract entity
	contractNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	// load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return err
	}

	workflowId := ""
	if contractEntity.PayOnline || contractEntity.PayAutomatically {
		workflowId = commonService.WorkflowInvoiceReadyNoPaymentLink
	} else {
		workflowId = commonService.WorkflowInvoiceReadyWithPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if contractEntity.PayOnline || contractEntity.PayAutomatically {
		paymentLink = h.cfg.Services.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := commonService.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: commonService.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(commonService.WorkflowInvoiceReadySubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
		},
		Attachments: []commonService.PostmarkEmailAttachment{},
	}

	err = h.appendInvoiceFileToEmailAsAttachment(ctx, eventData.Tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.services.CommonServices.PostmarkService.SendNotification(ctx, postmarkEmail, eventData.Tenant)

	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error sending invoice pay request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	h.createInvoiceAction(ctx, eventData.Tenant, invoiceEntity)

	// Request was successful
	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.SetPayInvoiceNotificationSentAt(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPayInvoiceNotificationSentAt"))
		h.log.Errorf("Error setting invoice pay notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) onInvoiceRemindNotificationV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.onInvoiceRemindNotificationV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceRemindNotificationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "evt.GetJsonData"))
		return errors.Wrap(err, "evt.GetJsonData")
	}
	tracing.LogObjectAsJson(span, "eventData", eventData)
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	//load invoice
	invoiceNode, err := h.services.CommonServices.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceEntity.IsOverdue() {
		tracing.TraceErr(span, errors.New("remind invoice notification requested for not applicable invoice"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		h.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract
	contractNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractForInvoice(ctx, eventData.Tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	//load tenant billing profile from neo4j
	tenantBillingProfileEntity, err := h.loadTenantBillingProfile(ctx, eventData.Tenant, false)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "loadTenantBillingProfile"))
		return err
	}
	tenantSettingsDbNode, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	workflowId := ""
	if invoiceEntity.PaymentDetails.PaymentLink == "" {
		workflowId = commonService.WorkflowInvoiceRemindNoPaymentLink
	} else {
		workflowId = commonService.WorkflowInvoiceRemindWithPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := utils.AddToListIfNotExists(contractEntity.InvoiceEmailBCC, tenantBillingProfileEntity.SendInvoicesBcc)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if invoiceEntity.PaymentDetails.PaymentLink != "" {
		paymentLink = h.cfg.Services.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := commonService.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: commonService.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(commonService.WorkflowInvoiceRemindSubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
		},
		Attachments: []commonService.PostmarkEmailAttachment{},
	}

	if tenantSettingsEntity.StripeCustomerPortalLink != "" {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = fmt.Sprintf(`PS: If you pay by card you can manage your billing details <a href="%s">here</a>.`, tenantSettingsEntity.StripeCustomerPortalLink)
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = `PS: If you pay by card you can manage your billing details here.`
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = tenantSettingsEntity.StripeCustomerPortalLink
	} else {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = ""
	}

	err = h.appendInvoiceFileToEmailAsAttachment(ctx, eventData.Tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendProviderLogoToEmail(ctx, eventData.Tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = h.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendCustomerOSLogoToEmail"))
		h.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = h.services.CommonServices.PostmarkService.SendNotification(ctx, postmarkEmail, eventData.Tenant)

	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		h.log.Errorf("Error sending invoice remind request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, eventData.Tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyLastRemindInvoiceNotificationSentAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
		h.log.Errorf("Error setting invoice remind notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) createInvoiceAction(ctx context.Context, tenant string, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("InvoiceId", invoiceEntity.Id))

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := neo4jenum.ActionInvoiceSent
	message := "Sent invoice N° " + invoiceEntity.Number + " with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)

	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ActionWriteRepository.MergeByActionType"))
		h.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (h *InvoiceEventHandler) appendInvoiceFileToEmailAsAttachment(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity, postmarkEmail *commonService.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendInvoiceFileToEmailAsAttachment")
	defer span.Finish()

	invoiceFileBytes, err := h.services.FileStoreApiService.GetFileBytes(tenant, invoice.RepositoryFileId, span)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, commonService.PostmarkEmailAttachment{
		Filename:       "Invoice " + invoice.Number + ".pdf",
		ContentEncoded: base64.StdEncoding.EncodeToString(*invoiceFileBytes),
		ContentType:    "application/pdf",
	})

	return nil
}

func (h *InvoiceEventHandler) appendProviderLogoToEmail(ctx context.Context, tenant, logoFileId string, postmarkEmail *commonService.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendProviderLogoToEmail")
	defer span.Finish()

	if logoFileId == "" {
		return nil
	}

	metadata, fileBytes, err := h.services.FileStoreApiService.GetFile(tenant, logoFileId, span)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, commonService.PostmarkEmailAttachment{
		Filename:       "provider-logo-file-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(*fileBytes),
		ContentType:    metadata.MimeType,
		ContentID:      "cid:provider-logo-file-encoded",
	})

	return nil
}

func (h *InvoiceEventHandler) appendCustomerOSLogoToEmail(ctx context.Context, postmarkEmail *commonService.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.appendCustomerOSLogoToEmail")
	defer span.Finish()

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd")
	}

	file, err := utils.GetFileByName(filepath.Join(currentDir, "/static", "customer-os.png"))

	b, err := ioutil.ReadAll(file)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ioutil.ReadAll"))
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, commonService.PostmarkEmailAttachment{
		Filename:       "customer-os-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(b),
		ContentType:    "image/png",
		ContentID:      "cid:customer-os-encoded",
	})

	return nil
}

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *InvoiceEventHandler) loadTenantBillingProfile(ctx context.Context, tenant string, failIfNotFound bool) (neo4jentity.TenantBillingProfileEntity, error) {
	tenantBillingProfiles, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, tenant)
	if err != nil {
		return neo4jentity.TenantBillingProfileEntity{}, err
	}
	if len(tenantBillingProfiles) == 0 {
		if failIfNotFound {
			return neo4jentity.TenantBillingProfileEntity{}, errors.New("tenantBillingProfiles not available")
		} else {
			return neo4jentity.TenantBillingProfileEntity{}, nil
		}
	}
	tenantBillingProfileEntity := neo4jmapper.MapDbNodeToTenantBillingProfileEntity(tenantBillingProfiles[0])
	return *tenantBillingProfileEntity, nil
}
