package listeners

import (
	"fmt"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neoEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/events-subscribers/listeners/website_visit_event"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

// add all subsribed agents for this handler here
var SubscribedAgents = [1]commonenum.AgentType{
	commonenum.AgentVisitorID,
}

var eventDataTypes = map[string]reflect.Type{
	data_fields.WebsiteVisitEvent{}.Type(): reflect.TypeOf(data_fields.WebsiteVisitEvent{}),
}

func OnWebhookEventCreated(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	_, webhookEvent, err := getWebhookEvent(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if webhookEvent == nil || webhookEvent.Data == nil || webhookEvent.DataType == "" {
		tracing.TraceErr(span, fmt.Errorf("invalid webhook event: missing required fields"))
		return fmt.Errorf("invalid webhook event: missing required fields")
	}

	// determine event handler //
	switch webhookEvent.DataType {

	// case data_fields.MeetingSummaryEvent{}.Type():
	// 	eventData, ok := webhookEvent.Data.(*data_fields.MeetingSummaryEvent)
	// 	if !ok {
	// 		return fmt.Errorf("failed to cast to MeetingSummaryEvent, got type: %T", webhookEvent.Data)
	// 	}
	//
	// 	err = handleMeetingSummaryEvent(ctx, dependencies.CommonServices, eventName, eventData)
	// 	if err != nil {
	// 		tracing.TraceErr(span, err)
	// 		return err
	// 	}
	case data_fields.WebsiteVisitEvent{}.Type():
		eventData, ok := webhookEvent.Data.(*data_fields.WebsiteVisitEvent)
		if !ok {
			tracing.TraceErr(span, fmt.Errorf("failed to cast to WebsiteVisitEvent, got type: %T", webhookEvent.Data))
			return fmt.Errorf("failed to cast to WebsiteVisitEvent, got type: %T", webhookEvent.Data)
		}

		websiteVisitEvent, err := website_visit_event.NewWebsiteVisitEventHandler(dependencies, eventData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return websiteVisitEvent.Handle(ctx)

	default:
		err = fmt.Errorf("Unsupported event %s", webhookEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
}

func handleMeetingSummaryEvent(ctx context.Context, s *service.CommonServices, sourceEvent commonenum.AgentListenerEvent, eventData *data_fields.MeetingSummaryEvent) error {
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

func publishCreateTimelineEvent(ctx context.Context, s *service.CommonServices, flow *postgres_entity.Flows, eventData *data_fields.MeetingSummaryEvent, sourceEvent commonenum.AgentListenerEvent) error {
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
	flowExecutionStatus, err := commonenum.GetFlowExecutionStatus(flowExecutionRecord.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowExecutionStatus != commonenum.FlowExecutionRunning {
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

func createFlowExecutionRecordForMeetingSummary(ctx context.Context, s *service.CommonServices, flow *postgres_entity.Flows, eventData *data_fields.MeetingSummaryEvent) (*postgres_entity.FlowExecution, error) {
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

func createMarkdownEvent(system commonenum.Source, eventData *data_fields.MeetingSummaryEvent) (data_fields.MarkdownEventFields, error) {
	var sourceId neoEntity.DataSource
	switch system {
	case commonenum.SourceFathom:
		sourceId = neoEntity.DataSourceFathom
	case commonenum.SourceGrain:
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
	system commonenum.Source, sourceEvent commonenum.AgentListenerEvent, orgId string,
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
		Name:             commonenum.AgentTimelineEventCreate,
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

func getWebhookEvent(input any) (commonenum.AgentListenerEvent, *dto.WebhookEvent, error) {
	// Cast input to Event
	message, ok := input.(*dto.Event)
	if !ok {
		return commonenum.NotSet, nil, fmt.Errorf("expected *dto.Event, got %T", input)
	}

	// Validate message
	if message == nil || message.Event.Data == nil {
		return commonenum.NotSet, nil, fmt.Errorf("message or message.Event.Data is nil")
	}

	// Cast Event.Data to WebhookEvent
	webhookEvent, ok := message.Event.Data.(*dto.WebhookEvent)
	if !ok {
		return commonenum.NotSet, nil, fmt.Errorf("expected *dto.WebhookEvent, got %T", message.Event.Data)
	}

	// Validate webhook event
	if webhookEvent.DataType == "" {
		return commonenum.NotSet, nil, fmt.Errorf("webhook event data type is empty")
	}

	// Get the expected type for this event
	eventDataType, ok := eventDataTypes[webhookEvent.DataType]
	if !ok {
		return commonenum.NotSet, nil, fmt.Errorf("unsupported event data type: %s", webhookEvent.DataType)
	}

	// Cast webhook data to map
	webhookData, ok := webhookEvent.Data.(map[string]interface{})
	if !ok {
		return commonenum.NotSet, nil, fmt.Errorf("expected map[string]interface{}, got %T", webhookEvent.Data)
	}

	// Create new instance of the target type
	webhookDataPtr := reflect.New(eventDataType).Interface()

	// Decode the map into the target type
	if err := utils.Decode(webhookData, webhookDataPtr); err != nil {
		return commonenum.NotSet, nil, fmt.Errorf("failed to decode webhook data: %w", err)
	}

	// Set the decoded data
	webhookEvent.Data = webhookDataPtr

	// Validate name before returning
	if webhookEvent.Name == commonenum.NotSet {
		return commonenum.NotSet, nil, fmt.Errorf("webhook event name is not set")
	}

	return webhookEvent.Name, webhookEvent, nil
}
