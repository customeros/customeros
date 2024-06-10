package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
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
	data := neo4jrepository.InteractionEventCreateFields{
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSourceOfTruth(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.Source),
			AppSource:     helper.GetAppSource(eventData.AppSource),
		},
		CreatedAt:          eventData.CreatedAt,
		Content:            eventData.Content,
		ContentType:        eventData.ContentType,
		Channel:            eventData.Channel,
		ChannelData:        eventData.ChannelData,
		Identifier:         eventData.Identifier,
		EventType:          eventData.EventType,
		BelongsToIssueId:   eventData.BelongsToIssueId,
		BelongsToSessionId: eventData.BelongsToSessionId,
		Hide:               eventData.Hide,
	}
	err := h.repositories.Neo4jRepositories.InteractionEventWriteRepository.Create(ctx, eventData.Tenant, interactionEventId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving interaction event %s: %s", interactionEventId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, interactionEventId, neo4jutil.NodeLabelInteractionEvent, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with external system %s: %s", interactionEventId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if eventData.Sender.Available() {
		err = h.repositories.Neo4jRepositories.InteractionEventWriteRepository.LinkInteractionEventWithSenderById(ctx, eventData.Tenant, interactionEventId, eventData.Sender.Participant.ID, eventData.Sender.Participant.NodeLabel(), eventData.Sender.RelationType)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link interaction event %s with sender %s: %s", interactionEventId, eventData.Sender.Participant.ID, err.Error())
		}
	}
	for _, receiver := range eventData.Receivers {
		if receiver.Available() {
			err = h.repositories.Neo4jRepositories.InteractionEventWriteRepository.LinkInteractionEventWithReceiverById(ctx, eventData.Tenant, interactionEventId, receiver.Participant.ID, receiver.Participant.NodeLabel(), receiver.RelationType)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while link interaction event %s with receiver %s: %s", interactionEventId, eventData.Sender.Participant.ID, err.Error())
			}
		}
	}

	organizationIds, orgsErr := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationIdsConnectedToInteractionEvent(ctx, eventData.Tenant, interactionEventId)
	if orgsErr != nil {
		tracing.TraceErr(span, orgsErr)
		h.log.Errorf("Error while getting organization ids connected to interaction event %s: %s", interactionEventId, orgsErr.Error())
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	for _, organizationId := range organizationIds {
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
				Tenant:         eventData.Tenant,
				OrganizationId: organizationId,
				AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			})
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
	data := neo4jrepository.InteractionEventUpdateFields{
		Content:     eventData.Content,
		ContentType: eventData.ContentType,
		Channel:     eventData.Channel,
		ChannelData: eventData.ChannelData,
		Identifier:  eventData.Identifier,
		EventType:   eventData.EventType,
		Hide:        eventData.Hide,
		Source:      helper.GetSourceOfTruth(eventData.Source),
	}
	err := h.repositories.Neo4jRepositories.InteractionEventWriteRepository.Update(ctx, eventData.Tenant, ieId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving interaction event %s: %s", ieId, err.Error())
	}

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, ieId, neo4jutil.NodeLabelInteractionEvent, externalSystemData)
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
	err := h.repositories.Neo4jRepositories.InteractionEventWriteRepository.SetAnalysisForInteractionEvent(ctx, eventData.Tenant, interactionEventId, eventData.Summary,
		eventData.ContentType, "summary", constants.SourceOpenline, constants.AppSourceEventProcessingPlatformSubscribers)
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
		err := h.repositories.Neo4jRepositories.InteractionEventWriteRepository.RemoveAllActionItemsForInteractionEvent(ctx, eventData.Tenant, interactionEventId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	var repoError error = nil
	for _, actionItem := range eventData.ActionItems {
		err := h.repositories.Neo4jRepositories.InteractionEventWriteRepository.AddActionItemForInteractionEvent(ctx, eventData.Tenant, interactionEventId, actionItem,
			constants.SourceOpenline, constants.AppSourceEventProcessingPlatformSubscribers)
		if err != nil {
			repoError = err
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while saving action items for email interaction event: %v", err)
		}
	}

	return repoError
}
