package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	orgcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphInteractionEventHandler struct {
	log                  logger.Logger
	organizationCommands *orgcmdhnd.OrganizationCommands
	repositories         *repository.Repositories
}

func (h *GraphInteractionEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnCreate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())

	var eventData event.InteractionEventCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	ieId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.InteractionEventRepository.Create(ctx, eventData.Tenant, ieId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving interaction event %s: %s", ieId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, ieId, constants.NodeLabel_InteractionEvent, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with external system %s: %s", ieId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	organizationIds, orgsErr := h.repositories.OrganizationRepository.GetOrganizationIdsConnectedToInteractionEvent(ctx, eventData.Tenant, ieId)
	if orgsErr != nil {
		tracing.TraceErr(span, orgsErr)
		h.log.Errorf("Error while getting organization ids connected to interaction event %s: %s", ieId, orgsErr.Error())
	}
	for _, organizationId := range organizationIds {
		err = h.organizationCommands.RefreshLastTouchpointCommand.Handle(ctx, cmd.NewRefreshLastTouchpointCommand(eventData.Tenant, organizationId, "", constants.AppSourceEventProcessingPlatform))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshLastTouchpointCommand failed: %v", err.Error())
		}
	}

	return nil
}

func (h *GraphInteractionEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnUpdate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())

	var eventData event.InteractionEventUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	ieId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.InteractionEventRepository.Update(ctx, eventData.Tenant, ieId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving interaction event %s: %s", ieId, err.Error())
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, ieId, constants.NodeLabel_InteractionEvent, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with external system %s: %s", ieId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return err
}

func (h *GraphInteractionEventHandler) OnSummaryReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnSummaryReplace")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.InteractionEventReplaceSummaryEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.InteractionEventRepository.SetAnalysisForInteractionEvent(ctx, eventData.Tenant, interactionEventId, eventData.Summary,
		eventData.ContentType, "summary", constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving analysis for email interaction event: %v", err)
	}

	return err
}

func (h *GraphInteractionEventHandler) OnActionItemsReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphInteractionEventHandler.OnActionItemsReplace")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.InteractionEventReplaceActionItemsEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.LogFields(log.Object("actionItems", eventData.ActionItems))

	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)

	if len(eventData.ActionItems) > 0 {
		err := h.repositories.InteractionEventRepository.RemoveAllActionItemsForInteractionEvent(ctx, eventData.Tenant, interactionEventId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	var repoError error = nil
	for _, actionItem := range eventData.ActionItems {
		err := h.repositories.InteractionEventRepository.AddActionItemForInteractionEvent(ctx, eventData.Tenant, interactionEventId, actionItem,
			constants.SourceOpenline, constants.AppSourceEventProcessingPlatform, eventData.UpdatedAt)
		if err != nil {
			repoError = err
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while saving action items for email interaction event: %v", err)
		}
	}

	return repoError
}
