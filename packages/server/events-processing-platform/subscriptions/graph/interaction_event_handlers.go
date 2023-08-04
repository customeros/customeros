package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphInteractionEventHandler struct {
	Log          logger.Logger
	Repositories *repository.Repositories
}

func (h *GraphInteractionEventHandler) OnSummaryReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnSummaryReplace")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.InteractionEventReplaceSummaryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.InteractionEventRepository.SetAnalysisForInteractionEvent(ctx, eventData.Tenant, interactionEventId, eventData.Summary,
		eventData.ContentType, "summary", constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.Log.Errorf("Error while saving analysis for email interaction event: %v", err)
	}

	return err
}

func (h *GraphInteractionEventHandler) OnActionItemsReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnActionItemsReplace")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.InteractionEventReplaceActionItemsEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.LogFields(log.Object("actionItems", eventData.ActionItems))

	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)

	if len(eventData.ActionItems) > 0 {
		err := h.Repositories.InteractionEventRepository.RemoveAllActionItemsForInteractionEvent(ctx, eventData.Tenant, interactionEventId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	var repoError error = nil
	for _, actionItem := range eventData.ActionItems {
		err := h.Repositories.InteractionEventRepository.AddActionItemForInteractionEvent(ctx, eventData.Tenant, interactionEventId, actionItem,
			constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, eventData.UpdatedAt)
		if err != nil {
			repoError = err
			tracing.TraceErr(span, err)
			h.Log.Errorf("Error while saving action items for email interaction event: %v", err)
		}
	}

	return repoError
}
