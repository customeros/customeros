package graph

import (
	"context"
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

func (h *InvoiceEventHandler) OnInvoiceNew(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceNew")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceNewEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	id := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.InvoiceNew(ctx, eventData.Tenant, eventData.ContractId, id, eventData.DryRun, eventData.Number, eventData.Date, eventData.DueDate, source, appSource, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving invoice %s: %s", id, err.Error())
		return err
	}
	return err
}

func (h *InvoiceEventHandler) OnInvoiceFill(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceFill")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	id := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.InvoiceFill(ctx, eventData.Tenant, id, eventData.Amount, eventData.VAT, eventData.Total, eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating invoice fill %s: %s", id, err.Error())
		return err
	}

	for _, item := range eventData.Lines {
		err := h.repositories.Neo4jRepositories.InvoiceWriteRepository.InvoiceFillInvoiceLine(ctx, eventData.Tenant, id, item.Index, item.Name, item.Price, item.Quantity, item.Amount, item.VAT, item.Total, eventData.UpdatedAt)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while inserting invoice line %s: %s", id, err.Error())
			return err
		}
	}

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
