package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

type InvoiceEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewInvoiceEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *InvoiceEventHandler) OnInvoiceCreateForContractV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceCreateForContractV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceForContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", eventData.ContractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	issuedDate := utils.ToDate(eventData.CreatedAt)
	if eventData.DryRun {
		issuedDate = utils.ToDate(eventData.PeriodStartDate)
		if eventData.Postpaid {
			issuedDate = utils.ToDate(eventData.PeriodEndDate).AddDate(0, 0, 1)
		}
	}
	dueDate := issuedDate.AddDate(0, 0, int(contractEntity.DueDays))
	data := neo4jrepository.InvoiceCreateFields{
		ContractId:           eventData.ContractId,
		Currency:             neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:               eventData.DryRun,
		OffCycle:             eventData.OffCycle,
		Postpaid:             eventData.Postpaid,
		Preview:              eventData.Preview,
		BillingCycleInMonths: eventData.BillingCycleInMonths,
		PeriodStartDate:      eventData.PeriodStartDate,
		PeriodEndDate:        eventData.PeriodEndDate,
		CreatedAt:            eventData.CreatedAt,
		IssuedDate:           issuedDate,
		DueDate:              dueDate,
		Status:               neo4jenum.InvoiceStatusInitialized,
		SourceFields: neo4jmodel.SourceFields{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		Note: eventData.Note,
	}
	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.CreateInvoiceForContract(ctx, eventData.Tenant, invoiceId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving invoice %s: %s", invoiceId, err.Error())
		return err
	}

	// Remove previous initialized invoices, if any
	if eventData.DryRun && eventData.Preview {
		err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeletePreviewCycleInitializedInvoices(ctx, eventData.Tenant, eventData.ContractId, invoiceId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while deleting preview invoice for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}
	}

	_ = h.callRequestFillInvoiceGRPC(ctx, eventData.Tenant, invoiceId, eventData.ContractId, span)

	return nil
}

func (h *InvoiceEventHandler) OnInvoiceFillV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceFillV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	invoiceEntityBeforeFill, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	data := neo4jrepository.InvoiceFillFields{
		ContractId:                   eventData.ContractId,
		Currency:                     neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:                       eventData.DryRun,
		Preview:                      eventData.Preview,
		OffCycle:                     eventData.OffCycle,
		Postpaid:                     eventData.Postpaid,
		InvoiceNumber:                eventData.InvoiceNumber,
		PeriodStartDate:              eventData.PeriodStartDate,
		PeriodEndDate:                eventData.PeriodEndDate,
		BillingCycleInMonths:         eventData.BillingCycleInMonths,
		Note:                         eventData.Note,
		CustomerName:                 eventData.Customer.Name,
		CustomerEmail:                eventData.Customer.Email,
		CustomerAddressLine1:         eventData.Customer.AddressLine1,
		CustomerAddressLine2:         eventData.Customer.AddressLine2,
		CustomerAddressZip:           eventData.Customer.Zip,
		CustomerAddressLocality:      eventData.Customer.Locality,
		CustomerAddressCountry:       eventData.Customer.Country,
		CustomerAddressRegion:        eventData.Customer.Region,
		ProviderLogoRepositoryFileId: eventData.Provider.LogoRepositoryFileId,
		ProviderName:                 eventData.Provider.Name,
		ProviderEmail:                eventData.Provider.Email,
		ProviderAddressLine1:         eventData.Provider.AddressLine1,
		ProviderAddressLine2:         eventData.Provider.AddressLine2,
		ProviderAddressZip:           eventData.Provider.Zip,
		ProviderAddressLocality:      eventData.Provider.Locality,
		ProviderAddressCountry:       eventData.Provider.Country,
		ProviderAddressRegion:        eventData.Provider.Region,
		Amount:                       eventData.Amount,
		VAT:                          eventData.VAT,
		TotalAmount:                  eventData.TotalAmount,
		Status:                       neo4jenum.DecodeInvoiceStatus(eventData.Status),
	}
	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.FillInvoice(ctx, eventData.Tenant, invoiceId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while filling invocie with details %s: %s", invoiceId, err.Error())
		return err
	}

	for _, item := range eventData.InvoiceLines {
		invoiceLineData := neo4jrepository.InvoiceLineCreateFields{
			CreatedAt:   item.CreatedAt,
			Name:        item.Name,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Amount:      item.Amount,
			VAT:         item.VAT,
			TotalAmount: item.TotalAmount,
			BilledType:  neo4jenum.DecodeBilledType(item.BilledType),
			SourceFields: neo4jmodel.SourceFields{
				Source:    helper.GetSource(item.SourceFields.Source),
				AppSource: helper.GetAppSource(item.SourceFields.AppSource),
			},
			ServiceLineItemId:       item.ServiceLineItemId,
			ServiceLineItemParentId: item.ServiceLineItemParentId,
		}
		err = h.services.CommonServices.Neo4jRepositories.InvoiceLineWriteRepository.CreateInvoiceLine(ctx, eventData.Tenant, invoiceId, item.Id, invoiceLineData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while inserting invoice line %s for invoice %s: %s", item.Id, invoiceId, err.Error())
			return err
		}
	}

	invoiceEntityAfterFill, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	if invoiceEntityAfterFill.Status != neo4jenum.InvoiceStatusEmpty {
		err = h.callGeneratePdfRequestGRPC(ctx, eventData.Tenant, invoiceId, span)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while calling generate pdf request for invoice %s: %s", invoiceId, err.Error())
			return err
		}
	}

	if !invoiceEntityAfterFill.OffCycle && !invoiceEntityAfterFill.DryRun {
		err := h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, eventData.Tenant, eventData.ContractId, "")
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while deleting preview invoice for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}

		err = h.callNextPreviewOnCycleInvoiceGRPC(ctx, eventData.Tenant, eventData.ContractId, span)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while calling next preview invoice for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}
	} else if invoiceEntityAfterFill.Preview && invoiceEntityAfterFill.DryRun {
		err := h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, eventData.Tenant, eventData.ContractId, invoiceId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while deleting preview invoice for contract %s: %s", eventData.ContractId, err.Error())
			return err
		}
	}

	if invoiceEntityAfterFill.Status != neo4jenum.InvoiceStatusEmpty {
		h.createInvoiceAction(ctx, eventData.Tenant, invoiceEntityBeforeFill.Status, *invoiceEntityAfterFill)
	}

	return nil
}

func (h *InvoiceEventHandler) callNextPreviewOnCycleInvoiceGRPC(ctx context.Context, tenant, contractId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return h.grpcClients.InvoiceClient.NextPreviewInvoiceForContract(ctx, &invoicepb.NextPreviewInvoiceForContractRequest{
			Tenant:     tenant,
			ContractId: contractId,
			AppSource:  constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error sending the next preview invoice request for contract {%s}: {%s}", contractId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) OnInvoicePdfGenerated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoicePdfGenerated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoicePdfGeneratedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	id := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, id)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	err := h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.InvoicePdfGenerated(ctx, eventData.Tenant, id, eventData.RepositoryFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating invoice pdf generated %s: %s", id, err.Error())
		return err
	}
	return err
}

func (h *InvoiceEventHandler) callGeneratePdfRequestGRPC(ctx context.Context, tenant, invoiceId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return h.grpcClients.InvoiceClient.GenerateInvoicePdf(ctx, &invoicepb.GenerateInvoicePdfRequest{
			Tenant:    tenant,
			InvoiceId: invoiceId,
			AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error sending the generate pdf request for invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) callRequestFillInvoiceGRPC(ctx context.Context, tenant, invoiceId, contractId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return h.grpcClients.InvoiceClient.RequestFillInvoice(ctx, &invoicepb.RequestFillInvoiceRequest{
			Tenant:     tenant,
			InvoiceId:  invoiceId,
			ContractId: contractId,
			AppSource:  constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error sending the request to fill invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return nil
}

func (h *InvoiceEventHandler) OnInvoiceVoidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceVoidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceVoidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	invoiceEntityBeforeVoid, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}

	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.UpdateInvoice(ctx, eventData.Tenant, invoiceId, neo4jrepository.InvoiceUpdateFields{
		UpdateStatus: true,
		Status:       neo4jenum.InvoiceStatusVoid,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while voiding invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}

	invoiceEntityAfterVoid, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	h.createInvoiceAction(ctx, eventData.Tenant, invoiceEntityBeforeVoid.Status, *invoiceEntityAfterVoid)

	return nil
}

func (h *InvoiceEventHandler) OnInvoiceDeleteV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceDeleteV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	err := h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeleteInitializedInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return err
}

func (h *InvoiceEventHandler) createInvoiceAction(ctx context.Context, tenant string, previousStatus neo4jenum.InvoiceStatus, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("invoiceId", invoiceEntity.Id))
	span.LogFields(log.String("previousStatus", previousStatus.String()))
	span.LogFields(log.String("newStatus", invoiceEntity.Status.String()))
	span.LogFields(log.Bool("dryRun", invoiceEntity.DryRun))
	span.LogFields(log.Float64("totalAmount", invoiceEntity.TotalAmount))

	if previousStatus == invoiceEntity.Status {
		return
	}
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		span.LogFields(log.String("result", "dry run or total amount is 0"))
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := neo4jenum.ActionNA
	message := ""
	switch invoiceEntity.Status {
	case neo4jenum.InvoiceStatusDue:
		message = "Invoice N° " + invoiceEntity.Number + " issued with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = neo4jenum.ActionInvoiceIssued
	case neo4jenum.InvoiceStatusPaid:
		message = "Invoice N° " + invoiceEntity.Number + " paid in full: " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = neo4jenum.ActionInvoicePaid
	case neo4jenum.InvoiceStatusVoid:
		message = "Invoice N° " + invoiceEntity.Number + " voided"
		actionType = neo4jenum.ActionInvoiceVoided
	case neo4jenum.InvoiceStatusOverdue:
		message = "Invoice N° " + invoiceEntity.Number + " overdue"
		actionType = neo4jenum.ActionInvoiceOverdue
	}
	if actionType == neo4jenum.ActionNA {
		span.LogFields(log.String("result", "status not supported"))
		return
	}
	if invoiceEntity.Status == neo4jenum.InvoiceStatusDue {
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	} else {
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (h *InvoiceEventHandler) getInvoice(ctx context.Context, tenant, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	invoiceDbNode, err := h.services.CommonServices.Neo4jRepositories.InvoiceReadRepository.GetInvoiceById(ctx, tenant, invoiceId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
}
