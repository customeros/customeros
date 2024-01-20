package graph

import (
	"context"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
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

func (h *InvoiceEventHandler) OnInvoiceCreateV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceCreateV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceCreateEvent
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
		err := h.repositories.Neo4jRepositories.InvoiceLineWriteRepository.CreateInvoiceLine(ctx, eventData.Tenant, id, item.Name, item.Price, item.Quantity, item.Amount, item.VAT, item.Total, eventData.UpdatedAt)
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
