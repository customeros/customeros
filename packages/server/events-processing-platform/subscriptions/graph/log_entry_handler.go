package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphLogEntryEventHandler struct {
	Log          logger.Logger
	Repositories *repository.Repositories
}

func (h *GraphLogEntryEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLogEntryEventHandler.OnCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.LogEntryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.LogEntryRepository.Create(ctx, eventData.Tenant, logEntryId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.Log.Errorf("Error while saving log entry %s: %s", logEntryId, err.Error())
	}

	return err
}

func (h *GraphLogEntryEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLogEntryEventHandler.OnCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.LogEntryUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.LogEntryRepository.Update(ctx, eventData.Tenant, logEntryId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.Log.Errorf("Error while saving log entry %s: %s", logEntryId, err.Error())
	}

	return err
}

func (h *GraphLogEntryEventHandler) OnAddTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLogEntryEventHandler.OnAddTag")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.LogEntryAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.TagRepository.AddTagByIdTo(ctx, eventData.Tenant, eventData.TagId, logEntryId, "LogEntry", eventData.TaggedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.Log.Errorf("Error while adding tag %s to log entry %s: %s", eventData.TagId, logEntryId, err.Error())
	}

	return err
}

func (h *GraphLogEntryEventHandler) OnRemoveTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLogEntryEventHandler.OnRemoveTag")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.LogEntryRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.TagRepository.RemoveTagByIdFrom(ctx, eventData.Tenant, eventData.TagId, logEntryId, "LogEntry")
	if err != nil {
		tracing.TraceErr(span, err)
		h.Log.Errorf("Error while removing tag %s to log entry %s: %s", eventData.TagId, logEntryId, err.Error())
	}

	return err
}
