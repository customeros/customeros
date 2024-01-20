package graph

import (
	"context"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type InvoiceEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewInvoiceEventHandler(log logger.Logger, repositories *repository.Repositories) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:          log,
		repositories: repositories,
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
		SourceFields: neo4jmodel.Source{
			Source:    eventData.SourceFields.Source,
			AppSource: eventData.SourceFields.AppSource,
		},
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
		Amount:          eventData.Amount,
		VAT:             eventData.VAT,
		TotalAmount:     eventData.TotalAmount,
		UpdatedAt:       eventData.UpdatedAt,
		ContractId:      eventData.ContractId,
		Currency:        neo4jenum.DecodeCurrency(eventData.Currency),
		DryRun:          eventData.DryRun,
		InvoiceNumber:   eventData.InvoiceNumber,
		PeriodStartDate: eventData.PeriodStartDate,
		PeriodEndDate:   eventData.PeriodEndDate,
		BillingCycle:    neo4jenum.DecodeBillingCycle(eventData.BillingCycle),
	}
	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.InvoiceFill(ctx, eventData.Tenant, invoiceId, data)
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

	// TODO remove technical flag from contract

	// TODO here generate request for invoice pdf generation

	return err
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
