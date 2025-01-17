package handlers

import (
	"context"
	"fmt"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	service "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neoEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

func HandleMeetingSummaryEvent(ctx context.Context, s *service.CommonServices, sourceEvent commonEnum.FlowListenerEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventData.Tenant,
	})

	flows, err := s.WorkflowService.GetFlowsForListenerEvent(ctx, sourceEvent, eventData.Type(), eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flows == nil || len(flows) == 0 {
		return nil
	}

	var errs error

	return errs
}

func publishCreateTimelineEvent(ctx context.Context, s *service.CommonServices, flow *entity.Flows, eventData *data_fields.MeetingSummaryEvent, sourceEvent commonEnum.FlowListenerEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.publishCreateTimelineEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// create contacts and org if they do not exist for all unique, non-tenant orgs
	allOrgIds, err := createContactsAndOrganizations(ctx, s, eventData.ParticipantEmails)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// build mdEvent
	source, err := sourceEvent.ExternalSystem()
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	mdEvent, err := createMarkdownEvent(source, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// create flow execution record
	flowExecutionRecord, err := createFlowExecutionRecordForMeetingSummary(ctx, s, flow, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// do not publish if flow is configured but not running
	flowExecutionStatus, err := commonEnum.GetFlowExecutionStatus(flowExecutionRecord.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowExecutionStatus != commonEnum.FlowExecutionRunning {
		return nil
	}

	// publish
	var allErr error
	for _, orgId := range *allOrgIds {
		err := handleMarkdownEventPublishing(ctx, s, source, sourceEvent, orgId, &mdEvent, flowExecutionRecord.ID)
		if err != nil {
			tracing.TraceErr(span, err)
			allErr = multierr.Append(allErr, err)
		}
	}

	return allErr
}

func createFlowExecutionRecordForMeetingSummary(ctx context.Context, s *service.CommonServices, flow *entity.Flows, eventData *data_fields.MeetingSummaryEvent) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.createFlowExecutionRecordForMeetingSummary")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	return nil, nil
}

func createContactsAndOrganizations(ctx context.Context, s *service.CommonServices, participantEmails *[]string) (*[]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.createContactsAndOrganizations")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	tenantDomains, err := s.WorkspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var allErr error
	allOrgIds := make([]string, 0)

	for _, email := range *participantEmails {

		// check that email does not belong to tenant
		vaildate := mailvalidate.ValidateEmailSyntax(email)
		if utils.IsStringInSlice(vaildate.Domain, tenantDomains) || !vaildate.IsValid {
			continue
		}

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

	return &allOrgIds, nil
}

func createOrgFromEmail(ctx context.Context, s *service.CommonServices, email string, tenantDomains []string) (string, error) {
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

func createMarkdownEvent(system commonEnum.Source, eventData *data_fields.MeetingSummaryEvent) (data_fields.MarkdownEventFields, error) {
	var sourceId neoEntity.DataSource
	switch system {
	case commonEnum.SourceFathom:
		sourceId = neoEntity.DataSourceFathom
	case commonEnum.SourceGrain:
		sourceId = neoEntity.DataSourceGrain
	default:
		return data_fields.MarkdownEventFields{}, fmt.Errorf("unsupported source: %v", system)
	}

	return data_fields.MarkdownEventFields{
		Source:    &sourceId,
		Content:   eventData.Content,
		CreatedAt: eventData.Timestamp,
	}, nil
}

func handleMarkdownEventPublishing(ctx context.Context, s *service.CommonServices,
	system commonEnum.Source, sourceEvent commonEnum.FlowListenerEvent, orgId string,
	mdEvent *data_fields.MarkdownEventFields, flowExecutionId string,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.handleMarkdownEventPublishing")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	mdEvent.OrganizationId = &orgId

	// build flow action event
	flowAgentEvent := dto.FlowAgentEvent{
		FlowExecutionId:  flowExecutionId,
		Tenant:           "",
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonEnum.AgentTimelineEventCreate,
		DataType:         data_fields.MarkdownEventFields{}.Type(),
		Data:             mdEvent,
	}

	// publish event
	if err := s.Events.Publisher.PublishFlowAgentEvent(ctx, flowAgentEvent); err != nil {
		err = fmt.Errorf("failed to publish markdown event for org %s: %w", orgId, err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
