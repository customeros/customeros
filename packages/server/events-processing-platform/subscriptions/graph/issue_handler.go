package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
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

type GraphIssueEventHandler struct {
	log                  logger.Logger
	organizationCommands *orgcmdhnd.OrganizationCommands
	Repositories         *repository.Repositories
}

func (h *GraphIssueEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphIssueEventHandler.OnCreate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())

	var eventData event.IssueCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.IssueRepository.Create(ctx, eventData.Tenant, issueId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.Repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, constants.NodeLabel_Issue, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link issue %s with external system %s: %s", issueId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if eventData.ReportedByOrganizationId != "" {
		err = h.organizationCommands.RefreshLastTouchpointCommand.Handle(ctx, cmd.NewRefreshLastTouchpointCommand(eventData.Tenant, eventData.ReportedByOrganizationId, "", constants.AppSourceEventProcessingPlatform))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshLastTouchpointCommand failed: %v", err.Error())
		}
	}

	return nil
}

func (h *GraphIssueEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphIssueEventHandler.OnUpdate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())

	var eventData event.IssueUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.LogFields(log.String("eventData", fmt.Sprintf("%+v", evt)))

	issueId := aggregate.GetIssueObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.IssueRepository.Update(ctx, eventData.Tenant, issueId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving issue %s: %s", issueId, err.Error())
	}

	if eventData.ExternalSystem.Available() {
		err = h.Repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, issueId, constants.NodeLabel_Issue, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link issue %s with external system %s: %s", issueId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return err
}
