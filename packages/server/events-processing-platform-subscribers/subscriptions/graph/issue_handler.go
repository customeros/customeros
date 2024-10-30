package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type IssueEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewIssueEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *IssueEventHandler {
	return &IssueEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

func (h *IssueEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.IssueCreateFields{
		CreatedAt: eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source),
			AppSource:     helper.GetAppSource(eventData.AppSource),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source),
		},
		GroupId:                   eventData.GroupId,
		Subject:                   eventData.Subject,
		Description:               eventData.Description,
		Status:                    eventData.Status,
		Priority:                  eventData.Priority,
		ReportedByOrganizationId:  eventData.ReportedByOrganizationId,
		SubmittedByOrganizationId: eventData.SubmittedByOrganizationId,
		SubmittedByUserId:         eventData.SubmittedByUserId,
	}
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.Create(ctx, eventData.Tenant, issueId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
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
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, model.NodeLabelIssue, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link issue %s with external system %s: %s", issueId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	if eventData.ReportedByOrganizationId != "" {
		err = h.services.CommonServices.OrganizationService.RequestRefreshLastTouchpoint(ctx, eventData.ReportedByOrganizationId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", eventData.ReportedByOrganizationId, err.Error())
		}
	}

	return nil
}

func (h *IssueEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.IssueUpdateFields{
		GroupId:     eventData.GroupId,
		Subject:     eventData.Subject,
		Description: eventData.Description,
		Status:      eventData.Status,
		Priority:    eventData.Priority,
		Source:      helper.GetSource(eventData.Source),
	}
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.Update(ctx, eventData.Tenant, issueId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
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
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, model.NodeLabelIssue, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link issue %s with external system %s: %s", issueId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return err
}

func (h *IssueEventHandler) OnAddUserAssignee(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnAddUserAssignee")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueAddUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.AddUserAssignee(ctx, eventData.Tenant, issueId, eventData.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding assignee to issue %s: %s", issueId, err.Error())
	}

	return err
}

func (h *IssueEventHandler) OnAddUserFollower(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnAddUserFollower")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueAddUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.AddUserFollower(ctx, eventData.Tenant, issueId, eventData.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adding follower to issue %s: %s", issueId, err.Error())
	}

	return err
}

func (h *IssueEventHandler) OnRemoveUserAssignee(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnRemoveUserAssignee")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueRemoveUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.RemoveUserAssignee(ctx, eventData.Tenant, issueId, eventData.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing assignee from issue %s: %s", issueId, err.Error())
	}

	return err

}

func (h *IssueEventHandler) OnRemoveUserFollower(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueEventHandler.OnRemoveUserFollower")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.IssueRemoveUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.RemoveUserFollower(ctx, eventData.Tenant, issueId, eventData.UserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing follower from issue %s: %s", issueId, err.Error())
	}

	return err
}
