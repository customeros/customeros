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
	flows, err := s.WorkflowService.GetFlowsByTrigger(ctx, sourceEvent)
	// send to dead events if no flows configured to receive event
	if err != nil || len(*flows) == 0 {
		err := sendToDeadEvents(ctx, s, sourceEvent, eventData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return err
	}

	var errs error
	for _, flow := range *flows {
		// get next action on the flow
		nextStep, err := s.WorkflowService.GetNextStepInFlow(ctx, flow.ID, flow.TriggerNodeID)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, fmt.Errorf("failed to get next step for flow %s: %w", flow.ID, err))
			continue
		}

		// validate the transition to next action
		validAction, err := s.WorkflowService.ValidateTransition(ctx, flow.TriggerNodeID, nextStep.ToNodeID)
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
		switch nextStep.ToNodeAction {
		case commonEnum.ActionTimelineEventCreate:
			if err := publishTimelineEventCreateEvent(ctx, s, flow.Status, flow.ID,
				nextStep.ToNodeID, eventData, sourceEvent); err != nil {
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

func publishTimelineEventCreateEvent(
	ctx context.Context, s *service.Services, flowStatus, flowId, flowNodeId string,
	eventData *data_fields.MeetingSummaryEvent, sourceEvent commonEnum.FlowListenerEvent,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.handleTimelineEventCreateAction")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	record, err := buildFlowExecutionRecord(ctx, flowStatus, flowId, flowNodeId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// save flow execution record
	flowExecutionRecord, err := s.WorkflowService.SaveFlowExecutionRecord(ctx, record)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionStatus, err := commonEnum.GetFlowExecutionStatus(flowExecutionRecord.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// if flow is configured but not running
	if executionStatus != commonEnum.FlowExecutionRunning {
		return nil
	}

	// process action execution event
	if err := processActionExecutionEvent(ctx, s, sourceEvent, flowExecutionRecord.ID, eventData); err != nil {
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
		Tenant:            common.GetTenantFromContext(ctx),
		FlowID:            flowId,
		EntityID:          eventData.MeetingID,
		EntityType:        commonEnum.EntityMeeting.String(),
		Status:            commonEnum.FlowExecutionRunning.String(),
		StartedAt:         utils.NowPtr(),
		CurrentStep:       commonEnum.ActionTimelineEventCreate.String(),
		CurrentStepNodeId: flowNodeId,
		CreatedAt:         utils.Now(),
		Context:           &data,
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

	_, err = s.PostgresRepositories.FlowDeadEventsRepository.Create(ctx, deadEvent)
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

func processActionExecutionEvent(ctx context.Context, s *service.Services, sourceEvent commonEnum.FlowListenerEvent, flowExecutionId string, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.processActionExecutionEvent")
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
	allOrgIds := make([]string, 0)

	for _, email := range *eventData.ParticipantEmails {

		// create the org if if doesn't exist
		organizationId, err := createOrgFromEmail(ctx, s, email, tenantDomains)
		if err != nil {
			tracing.TraceErr(span, err)
			allErr = multierr.Append(allErr, err)
		}

		// if org is unique, add to all org slice
		if organizationId != "" && !utils.IsStringInSlice(organizationId, allOrgIds) {
			allOrgIds = append(allOrgIds, organizationId)
		}

		// create the contact if it does not exist
		_, err = s.ContactService.CreateContactWithOrganizationByEmail(ctx, nil, email)
		if err != nil {
			err = fmt.Errorf("failed to create contact from email %s: %w", email, err)
			tracing.TraceErr(span, err)
			allErr = multierr.Append(allErr, err)
		}
	}

	// publish event for all unique non-tenant orgs
	for _, orgId := range allOrgIds {
		err := handleMarkdownEventPublishing(ctx, s, system, sourceEvent, orgId, &mdEvent, flowExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			allErr = multierr.Append(allErr, err)
		}
	}

	return allErr
}

func createOrgFromEmail(ctx context.Context, s *service.Services, email string, tenantDomains []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.createOrgFromEmail")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	cleanEmail := mailvalidate.ValidateEmailSyntax(email)
	if !cleanEmail.IsValid {
		err := fmt.Errorf("invalid email address: %s", email)
		tracing.TraceErr(span, err)
		return "", err
	}

	if utils.IsStringInSlice(cleanEmail.Domain, tenantDomains) {
		return "", nil
	}

	id, err := s.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{cleanEmail.Domain},
	})
	if err != nil {
		err = fmt.Errorf("failed to save organization for domain %s: %w", cleanEmail.Domain, err)
		tracing.TraceErr(span, err)
		return "", err
	}

	return id, nil
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

func handleMarkdownEventPublishing(ctx context.Context, s *service.Services,
	system enum.ExternalSystemId, sourceEvent commonEnum.FlowListenerEvent, orgId string,
	mdEvent *data_fields.MarkdownEventFields, flowExecutionId string,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.handleMarkdownEventPublishing")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	mdEvent.OrganizationId = &orgId

	// build flow action event
	flowActionEvent := dto.FlowActionEvent{
		FlowExecutionId:  flowExecutionId,
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonEnum.ActionTimelineEventCreate,
		DataType:         data_fields.MarkdownEventFields{}.Type(),
		Data:             mdEvent,
	}

	// publish event
	if err := s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent); err != nil {
		err = fmt.Errorf("failed to publish markdown event for org %s: %w", orgId, err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
