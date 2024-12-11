package handlers

import (
	"context"
	"fmt"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func HandleMeetingSummaryEvent(ctx context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// check to see if tenant has flows configured for this event
	flows, err := s.WorkflowService.GetWorkflowsByListenerEvent(ctx, sourceEvent)
	// send to dead events if no flows configured to receive event
	if err != nil || len(flows) == 0 {
		err := sendToDeadEvents(ctx, s, sourceEvent, eventData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	var errs error
	for _, flow := range flows {
		// get next action on the flow
		nextFlowAction, nextFlowNodeId, err := s.WorkflowService.GetFirstAction(ctx, &flow)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, fmt.Errorf("failed to get first action for flow %s: %w", flow.ID, err))
			continue
		}

		// validate the transition to next action
		validAction, err := s.WorkflowService.IsFlowActionValidTransitionFromListener(ctx, sourceEvent, nextFlowAction)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, fmt.Errorf("failed to validate transition for flow %s: %w", flow.ID, err))
			continue
		}
		if !validAction {
			err = fmt.Errorf("not a valid action transition for flow %s", flow.ID)
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
			continue
		}

		// prepare and send event to action executioner
		switch nextFlowAction {
		case commonEnum.ActionTimelineEventCreate:
			if err := sendTimelineEventCreateAction(ctx, s, flow.Status.String(), flow.ID,
				nextFlowNodeId, eventData, sourceEvent); err != nil {
				tracing.TraceErr(span, err)
				errs = multierr.Append(errs, fmt.Errorf("failed to send timeline event for flow %s: %w", flow.ID, err))
			}
		default:
			err = fmt.Errorf("next flow action not handled for flow %s", flow.ID)
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func sendTimelineEventCreateAction(
	ctx context.Context, s *service.Services, flowStatus, flowId, flowNodeId string,
	eventData *data_fields.MeetingSummaryEvent, sourceEvent commonEnum.FlowListenerEvent,
) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.hantdleTimelineEventCreateAction")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	record, err := buildFlowExecutionRecord(ctx, flowStatus, flowId, flowNodeId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// save flow execution record
	flowExecutionID, err := s.PostgresRepositories.FlowExecutionRepository.Save(ctx, record)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionStatus, err := commonEnum.GetFlowExecutionStatus(record.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if executionStatus != commonEnum.FlowExecutionRunning {
		return nil
	}

	// if not blocked, send publish action execution event
	if err := publishCreateMarkdownEvent(ctx, s, sourceEvent, flowExecutionID, eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish markdown event"))
		return err
	}

	return nil
}

func buildFlowExecutionRecord(ctx context.Context, flowStatus, flowId, flowNodeId string, eventData *data_fields.MeetingSummaryEvent) (postgresEntity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.buildFlowExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	data, err := eventData.ToString()
	if err != nil {
		tracing.TraceErr(span, err)
		return postgresEntity.FlowExecution{}, err
	}

	record := postgresEntity.FlowExecution{
		Tenant:         common.GetTenantFromContext(ctx),
		FlowID:         flowId,
		EntityID:       eventData.MeetingID,
		EntityType:     commonEnum.EntityMeeting.String(),
		Status:         commonEnum.FlowExecutionRunning.String(),
		StartedAt:      utils.NowPtr(),
		NextStep:       commonEnum.ActionTimelineEventCreate.String(),
		NextStepNodeId: flowNodeId,
		CreatedAt:      utils.Now(),
		Context:        &data,
	}

	switch flowStatus {
	case "INACTIVE":
		reason := commonEnum.FlowBlockedNotActive.String()
		record.Status = commonEnum.FlowExecutionBlocked.String()
		record.BlockedReason = &reason

	case "ARCHIVED":
		reason := commonEnum.FlowBlockedArchived.String()
		record.Status = commonEnum.FlowExecutionBlocked.String()
		record.BlockedReason = &reason

	default:
		record.Status = commonEnum.FlowExecutionRunning.String()
	}
	return record, nil
}

func sendToDeadEvents(ctx context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.sendToDeadEvents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	deadEvent, err := buildDeadEventFromMeetingSummary(ctx, sourceEvent, eventData)
	if err != nil {
		return err
	}

	_, err = s.PostgresRepositories.FlowDeadEventsRepository.Save(ctx, deadEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func buildDeadEventFromMeetingSummary(ctx context.Context, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.MeetingSummaryEvent) (postgresEntity.FlowDeadEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.buildDeadEventFromMeetingSummary")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	data, err := eventData.ToString()
	if err != nil {
		tracing.TraceErr(span, err)
		return postgresEntity.FlowDeadEvents{}, err
	}

	return postgresEntity.FlowDeadEvents{
		Tenant:    common.GetTenantFromContext(ctx),
		NodeType:  commonEnum.NodeFlowListenerEvent.String(),
		EventType: eventData.Type(),
		Event:     sourceEvent.String(),
		CreatedAt: utils.Now(),
		Data:      &data,
	}, nil
}

func publishCreateMarkdownEvent(ctx context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, flowExecutionId string, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.publishCreateMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	system, err := sourceEvent.ExternalSystem()
	if err != nil {
		return err
	}

	tenantDomains, err := s.WorkspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	mdEvent, err := createMarkdownEvent(system, eventData)
	if err != nil {
		return err
	}

	var allErr error
	for _, email := range *eventData.ParticipantEmails {
		if err := handleMarkdownEventPublishing(ctx, span, s, system, sourceEvent, email, tenantDomains, &mdEvent, flowExecutionId); err != nil {
			allErr = multierr.Append(allErr, err)
		}
	}

	return allErr
}

func createMarkdownEvent(system enum.ExternalSystemId, eventData *data_fields.MeetingSummaryEvent) (data_fields.MarkdownEventFields, error) {
	var sourceId entity.DataSource
	switch system {
	case enum.Fathom:
		sourceId = entity.DataSourceFathom
	case enum.Grain:
		sourceId = entity.DataSourceGrain
	default:
		return data_fields.MarkdownEventFields{}, fmt.Errorf("unsupported source: %v", system)
	}

	return data_fields.MarkdownEventFields{
		Source:    &sourceId,
		Content:   eventData.Content,
		CreatedAt: eventData.Timestamp,
	}, nil
}

func handleMarkdownEventPublishing(ctx context.Context, span opentracing.Span, s *service.Services,
	system enum.ExternalSystemId, sourceEvent commonEnum.FlowListenerEvent, email string, tenantDomains []string,
	mdEvent *data_fields.MarkdownEventFields, flowExecutionId string,
) error {
	cleanEmail := mailvalidate.ValidateEmailSyntax(email)
	if !cleanEmail.IsValid {
		return fmt.Errorf("invalid email address: %s", email)
	}

	if isEmailInDomains(cleanEmail.Domain, tenantDomains) {
		return nil
	}

	id, err := s.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{cleanEmail.Domain},
	})
	if err != nil {
		return fmt.Errorf("failed to save organization for domain %s: %w", cleanEmail.Domain, err)
	}

	mdEvent.OrganizationId = &id

	flowActionEvent := dto.FlowActionEvent{
		FlowExecutionId:  flowExecutionId,
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonEnum.ActionTimelineEventCreate,
		DataType:         data_fields.MarkdownEventFields{}.Type(),
		Data:             mdEvent,
	}

	if err := s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent); err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to publish markdown event for email %s: %w", email, err)
	}

	return nil
}

func isEmailInDomains(emailDomain string, tenantDomains []string) bool {
	for _, domain := range tenantDomains {
		if domain == emailDomain {
			return true
		}
	}
	return false
}

func createFlowActionEventCreateContact(system enum.ExternalSystemId, sourceEvent commonEnum.FlowListenerEvent, email string) dto.FlowActionEvent {
	return dto.FlowActionEvent{
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonEnum.ActionContactCreate,
		DataType:         data_fields.ContactCreateEvent{}.Type(),
		Data: data_fields.ContactCreateEvent{
			Email: email,
		},
	}
}
