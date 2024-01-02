package graph

import (
	"context"
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type IssueEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewIssueEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *IssueEventHandler {
	return &IssueEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
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
		UpdatedAt: eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source),
			AppSource:     helper.GetAppSource(eventData.AppSource),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source),
		},
		Subject:                   eventData.Subject,
		Description:               eventData.Description,
		Status:                    eventData.Status,
		Priority:                  eventData.Priority,
		ReportedByOrganizationId:  eventData.ReportedByOrganizationId,
		SubmittedByOrganizationId: eventData.SubmittedByOrganizationId,
		SubmittedByUserId:         eventData.SubmittedByUserId,
	}
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.Create(ctx, eventData.Tenant, issueId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, neo4jentity.NodeLabel_Issue, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link issue %s with external system %s: %s", issueId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	if eventData.ReportedByOrganizationId != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: eventData.ReportedByOrganizationId,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
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
		Subject:     eventData.Subject,
		Description: eventData.Description,
		Status:      eventData.Status,
		Priority:    eventData.Priority,
		UpdatedAt:   eventData.UpdatedAt,
		Source:      helper.GetSource(eventData.Source),
	}
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.Update(ctx, eventData.Tenant, issueId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, neo4jentity.NodeLabel_Issue, eventData.ExternalSystem)
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
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.AddUserAssignee(ctx, eventData.Tenant, issueId, eventData.UserId, eventData.At)
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
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.AddUserFollower(ctx, eventData.Tenant, issueId, eventData.UserId, eventData.At)
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
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.RemoveUserAssignee(ctx, eventData.Tenant, issueId, eventData.UserId, eventData.At)
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
	err := h.repositories.Neo4jRepositories.IssueWriteRepository.RemoveUserFollower(ctx, eventData.Tenant, issueId, eventData.UserId, eventData.At)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while removing follower from issue %s: %s", issueId, err.Error())
	}

	return err
}
