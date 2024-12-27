package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
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
