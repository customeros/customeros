package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	invoicingcycle "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type InvoicingCycleEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewInvoicingCycleEventHandler(log logger.Logger, repositories *repository.Repositories) *InvoicingCycleEventHandler {
	return &InvoicingCycleEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *InvoicingCycleEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoicingcycle.InvoicingCycleCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoicingCycleId := invoicingcycle.GetInvoicingCycleObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoicingCycleId)

	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.repositories.Neo4jRepositories.InvoicingCycleWriteRepository.CreateInvoicingCycleType(ctx, eventData.Tenant, invoicingCycleId, eventData.Type, source, appSource, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving invoicing cycle type %s: %s", invoicingCycleId, err.Error())
		return err
	}
	return err
}

func (h *InvoicingCycleEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicingCycleEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoicingcycle.InvoicingCycleUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	invoicingCycleId := invoicingcycle.GetInvoicingCycleObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoicingCycleId)

	err := h.repositories.Neo4jRepositories.InvoicingCycleWriteRepository.UpdateInvoicingCycleType(ctx, eventData.Tenant, invoicingCycleId, eventData.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating invoicing cycle type %s: %s", invoicingCycleId, err.Error())
		return err
	}
	return err
}
