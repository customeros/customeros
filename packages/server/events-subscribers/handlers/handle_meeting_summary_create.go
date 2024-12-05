package handlers

import (
	"context"
	"fmt"

	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func HandleMeetingSummaryEvent(ctx context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// TODO: lookup which actions are configured as part of flow for tenant
	// only trigger events for actions that are turned on
	if err := publishCreateMarkdownEvent(ctx, s, sourceEvent, eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish markdown event"))
		return err
	}

	if err := publishEventCreateContact(ctx, s, sourceEvent, eventData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish contact creation event"))
		return err
	}

	return nil
}

func publishEventCreateContact(ctx context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.publishEventCreateContact")
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

	var allErr error
	for _, email := range *eventData.ParticipantEmails {
		if err := handleContactEventPublishing(ctx, span, s, system, sourceEvent, email, tenantDomains); err != nil {
			allErr = multierr.Append(allErr, err)
		}
	}
	return allErr
}

func handleContactEventPublishing(ctx context.Context, span opentracing.Span, s *service.Services,
	system enum.ExternalSystemId, sourceEvent commonenum.FlowEvent, email string, tenantDomains []string,
) error {
	cleanEmail := mailvalidate.ValidateEmailSyntax(email)
	if !cleanEmail.IsValid {
		return fmt.Errorf("invalid email address: %s", email)
	}

	if isEmailInDomains(cleanEmail.Domain, tenantDomains) {
		return nil
	}

	flowActionEvent := createFlowActionEventCreateContact(system, sourceEvent, email)
	if err := s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent); err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to publish contact creation for email %s: %w", email, err)
	}

	return nil
}

func publishCreateMarkdownEvent(ctx context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
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
		if err := handleMarkdownEventPublishing(ctx, span, s, system, sourceEvent, email, tenantDomains, &mdEvent); err != nil {
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
	system enum.ExternalSystemId, sourceEvent commonenum.FlowEvent, email string, tenantDomains []string,
	mdEvent *data_fields.MarkdownEventFields,
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
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonenum.ActionTimelineEventCreate,
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

func createFlowActionEventCreateContact(system enum.ExternalSystemId, sourceEvent commonenum.FlowEvent, email string) dto.FlowActionEvent {
	return dto.FlowActionEvent{
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonenum.ActionContactCreate,
		DataType:         data_fields.ContactCreateEvent{}.Type(),
		Data: data_fields.ContactCreateEvent{
			Email: email,
		},
	}
}
