package graph

import (
	"context"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type InvoiceEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewInvoiceEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
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

	data := neo4jrepository.InvoiceCreateFields{
		ContractId:      eventData.ContractId,
		Currency:        neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:          eventData.DryRun,
		OffCycle:        eventData.OffCycle,
		Postpaid:        eventData.Postpaid,
		BillingCycle:    neo4jenum.DecodeBillingCycle(eventData.BillingCycle),
		PeriodStartDate: eventData.PeriodStartDate,
		PeriodEndDate:   eventData.PeriodEndDate,
		CreatedAt:       eventData.CreatedAt,
		DueDate:         eventData.DueDate,
		Status:          neo4jenum.InvoiceStatusDraft,
		SourceFields: neo4jmodel.Source{
			Source:    eventData.SourceFields.Source,
			AppSource: eventData.SourceFields.AppSource,
		},
		Note: eventData.Note,
	}
	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.CreateInvoiceForContract(ctx, eventData.Tenant, invoiceId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving invoice %s: %s", invoiceId, err.Error())
		return err
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

	data := neo4jrepository.InvoiceFillFields{
		UpdatedAt:                    eventData.UpdatedAt,
		ContractId:                   eventData.ContractId,
		Currency:                     neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:                       eventData.DryRun,
		InvoiceNumber:                eventData.InvoiceNumber,
		PeriodStartDate:              eventData.PeriodStartDate,
		PeriodEndDate:                eventData.PeriodEndDate,
		BillingCycle:                 neo4jenum.DecodeBillingCycle(eventData.BillingCycle),
		Note:                         eventData.Note,
		CustomerName:                 eventData.Customer.Name,
		CustomerEmail:                eventData.Customer.Email,
		CustomerAddressLine1:         eventData.Customer.AddressLine1,
		CustomerAddressLine2:         eventData.Customer.AddressLine2,
		CustomerAddressZip:           eventData.Customer.Zip,
		CustomerAddressLocality:      eventData.Customer.Locality,
		CustomerAddressCountry:       eventData.Customer.Country,
		ProviderLogoRepositoryFileId: eventData.Provider.LogoRepositoryFileId,
		ProviderName:                 eventData.Provider.Name,
		ProviderEmail:                eventData.Provider.Email,
		ProviderAddressLine1:         eventData.Provider.AddressLine1,
		ProviderAddressLine2:         eventData.Provider.AddressLine2,
		ProviderAddressZip:           eventData.Provider.Zip,
		ProviderAddressLocality:      eventData.Provider.Locality,
		ProviderAddressCountry:       eventData.Provider.Country,
		Amount:                       eventData.Amount,
		VAT:                          eventData.VAT,
		TotalAmount:                  eventData.TotalAmount,
		Status:                       neo4jenum.DecodeInvoiceStatus(eventData.Status),
	}
	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.FillInvoice(ctx, eventData.Tenant, invoiceId, data)
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
			SourceFields: neo4jmodel.Source{
				Source:    helper.GetSource(item.SourceFields.Source),
				AppSource: helper.GetAppSource(item.SourceFields.AppSource),
			},
			ServiceLineItemId:       item.ServiceLineItemId,
			ServiceLineItemParentId: item.ServiceLineItemParentId,
		}
		err = h.repositories.Neo4jRepositories.InvoiceLineWriteRepository.CreateInvoiceLine(ctx, eventData.Tenant, invoiceId, item.Id, invoiceLineData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while inserting invoice line %s for invoice %s: %s", item.Id, invoiceId, err.Error())
			return err
		}
	}

	err = h.callGeneratePdfRequestGRPC(ctx, eventData.Tenant, invoiceId, span)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while calling generate pdf request for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (h *InvoiceEventHandler) OnInvoiceUpdateV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceUpdateV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)

	data := neo4jrepository.InvoiceUpdateFields{
		UpdatedAt:         eventData.UpdatedAt,
		Status:            neo4jenum.DecodeInvoiceStatus(eventData.Status),
		PaymentLink:       eventData.PaymentLink,
		UpdateStatus:      eventData.UpdateStatus(),
		UpdatePaymentLink: eventData.UpdatePaymentLink(),
	}
	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.UpdateInvoice(ctx, eventData.Tenant, invoiceId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating invoice %s: %s", invoiceId, err.Error())
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

	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.InvoicePdfGenerated(ctx, eventData.Tenant, id, eventData.RepositoryFileId, eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating invoice pdf generated %s: %s", id, err.Error())
		return err
	}
	return err
}

func (s *InvoiceEventHandler) callGeneratePdfRequestGRPC(ctx context.Context, tenant, invoiceId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.grpcClients.InvoiceClient.GenerateInvoicePdf(ctx, &invoicepb.GenerateInvoicePdfRequest{
			Tenant:    tenant,
			InvoiceId: invoiceId,
			AppSource: constants.AppSourceEventProcessingPlatform,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error sending the generate pdf request for invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return nil
}

func (s *InvoiceEventHandler) callRequestFillInvoiceGRPC(ctx context.Context, tenant, invoiceId, contractId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*invoicepb.InvoiceIdResponse](func() (*invoicepb.InvoiceIdResponse, error) {
		return s.grpcClients.InvoiceClient.RequestFillInvoice(ctx, &invoicepb.RequestFillInvoiceRequest{
			Tenant:     tenant,
			InvoiceId:  invoiceId,
			ContractId: contractId,
			AppSource:  constants.AppSourceEventProcessingPlatform,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error sending the request to fill invoice {%s}: {%s}", invoiceId, err.Error())
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

	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.UpdateInvoice(ctx, eventData.Tenant, invoiceId, neo4jrepository.InvoiceUpdateFields{
		UpdatedAt:    eventData.UpdatedAt,
		UpdateStatus: true,
		Status:       neo4jenum.InvoiceStatusVoid,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while voiding invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return err
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

	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.DeleteInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return err
}
