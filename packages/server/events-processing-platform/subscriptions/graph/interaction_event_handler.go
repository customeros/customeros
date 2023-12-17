package graph

import (
	"context"
	"fmt"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InteractionEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewInteractionEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *InteractionEventHandler {
	return &InteractionEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (h *InteractionEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.InteractionEventCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	interactionEventId := aggregate.GetInteractionEventObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.InteractionEventRepository.Create(ctx, eventData.Tenant, interactionEventId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving interaction event %s: %s", interactionEventId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, interactionEventId, constants.NodeLabel_InteractionEvent, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with external system %s: %s", interactionEventId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if eventData.Sender.Available() {
		err = h.repositories.InteractionEventRepository.LinkInteractionEventWithSenderById(ctx, eventData.Tenant, interactionEventId, eventData.Sender.Participant.ID, eventData.Sender.Participant.NodeLabel(), eventData.Sender.RelationType)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with sender %s: %s", interactionEventId, eventData.Sender.Participant.ID, err.Error())
		}
	}
	for _, receiver := range eventData.Receivers {
		if receiver.Available() {
			err = h.repositories.InteractionEventRepository.LinkInteractionEventWithReceiverById(ctx, eventData.Tenant, interactionEventId, receiver.Participant.ID, receiver.Participant.NodeLabel(), receiver.RelationType)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while link interaction event %s with receiver %s: %s", interactionEventId, eventData.Sender.Participant.ID, err.Error())
			}
		}
	}

	organizationIds, orgsErr := h.repositories.OrganizationRepository.GetOrganizationIdsConnectedToInteractionEvent(ctx, eventData.Tenant, interactionEventId)
	if orgsErr != nil {
		tracing.TraceErr(span, orgsErr)
		h.log.Errorf("Error while getting organization ids connected to interaction event %s: %s", interactionEventId, orgsErr.Error())
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	for _, organizationId := range organizationIds {
		_, err = h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organizationId,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", organizationId, err.Error())
		}
	}

	return nil
}

func (h *InteractionEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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

func (h *InteractionEventHandler) OnSummaryReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.OnSummaryReplace")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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

func (h *InteractionEventHandler) OnActionItemsReplace(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventHandler.OnActionItemsReplace")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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
