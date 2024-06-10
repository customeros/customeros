package graph

import (
	"context"
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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type LogEntryEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewLogEntryEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *LogEntryEventHandler {
	return &LogEntryEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (h *LogEntryEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.LogEntryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.LogEntryCreateFields{
		Content:              eventData.Content,
		ContentType:          eventData.ContentType,
		StartedAt:            eventData.StartedAt,
		AuthorUserId:         eventData.AuthorUserId,
		LoggedOrganizationId: eventData.LoggedOrganizationId,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.AppSource),
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.repositories.Neo4jRepositories.LogEntryWriteRepository.Create(ctx, eventData.Tenant, logEntryId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving log entry %s: %s", logEntryId, err.Error())
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
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, logEntryId, neo4jutil.NodeLabelLogEntry, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link log entry %s with external system %s: %s", logEntryId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: eventData.LoggedOrganizationId,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", eventData.LoggedOrganizationId, err.Error())
	}

	return nil
}

func (h *LogEntryEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.LogEntryUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.LogEntryUpdateFields{
		Content:              eventData.Content,
		ContentType:          eventData.ContentType,
		StartedAt:            eventData.StartedAt,
		LoggedOrganizationId: eventData.LoggedOrganizationId,
		Source:               helper.GetSource(eventData.SourceOfTruth),
	}
	err := h.repositories.Neo4jRepositories.LogEntryWriteRepository.Update(ctx, eventData.Tenant, logEntryId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving log entry %s: %s", logEntryId, err.Error())
	}

	return err
}

func (h *LogEntryEventHandler) OnAddTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryEventHandler.OnAddTag")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.LogEntryAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.TagWriteRepository.LinkTagByIdToEntity(ctx, eventData.Tenant, eventData.TagId, logEntryId, "LogEntry", eventData.TaggedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding tag %s to log entry %s: %s", eventData.TagId, logEntryId, err.Error())
	}

	return err
}

func (h *LogEntryEventHandler) OnRemoveTag(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryEventHandler.OnRemoveTag")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.LogEntryRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	logEntryId := aggregate.GetLogEntryObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.TagWriteRepository.UnlinkTagByIdFromEntity(ctx, eventData.Tenant, eventData.TagId, logEntryId, "LogEntry")
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing tag %s to log entry %s: %s", eventData.TagId, logEntryId, err.Error())
	}

	return err
}
