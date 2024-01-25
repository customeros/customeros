package graph

import (
	"context"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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
		InvoiceNumber:   eventData.InvoiceNumber,
		BillingCycle:    neo4jenum.DecodeBillingCycle(eventData.BillingCycle),
		PeriodStartDate: eventData.PeriodStartDate,
		PeriodEndDate:   eventData.PeriodEndDate,
		CreatedAt:       eventData.CreatedAt,
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
	return err
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
		UpdatedAt:                     eventData.UpdatedAt,
		ContractId:                    eventData.ContractId,
		Currency:                      neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:                        eventData.DryRun,
		InvoiceNumber:                 eventData.InvoiceNumber,
		PeriodStartDate:               eventData.PeriodStartDate,
		PeriodEndDate:                 eventData.PeriodEndDate,
		BillingCycle:                  neo4jenum.DecodeBillingCycle(eventData.BillingCycle),
		Note:                          eventData.Note,
		DomesticPaymentsBankInfo:      eventData.DomesticPaymentsBankInfo,
		InternationalPaymentsBankInfo: eventData.InternationalPaymentsBankInfo,
		CustomerName:                  eventData.Customer.Name,
		CustomerEmail:                 eventData.Customer.Email,
		CustomerAddress:               eventData.Customer.Address,
		ProviderLogoUrl:               eventData.Provider.LogoUrl,
		ProviderName:                  eventData.Provider.Name,
		ProviderAddress:               eventData.Provider.Address,
		Amount:                        eventData.Amount,
		VAT:                           eventData.VAT,
		TotalAmount:                   eventData.TotalAmount,
		Status:                        neo4jenum.DecodeInvoiceStatus(eventData.Status),
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
		UpdatedAt:    eventData.UpdatedAt,
		Status:       neo4jenum.DecodeInvoiceStatus(eventData.Status),
		UpdateStatus: eventData.UpdateStatus(),
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
	_, err := s.grpcClients.InvoiceClient.GenerateInvoicePdf(ctx, &invoicepb.GenerateInvoicePdfRequest{
		Tenant:    tenant,
		InvoiceId: invoiceId,
		AppSource: constants.AppSourceEventProcessingPlatform,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error sending the generate pdf request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}
