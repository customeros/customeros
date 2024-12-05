package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

func HandleMeetingSummaryEvent(ctx context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// todo lookup which actions are configured as part of flow for tenant
	// only trigger events for actions that are turned on

	err := publishCreateMarkdownEvent(ctx, s, sourceEvent, eventData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish markdown event"))
		return err
	}

	err = publishEventCreateContact(ctx, s, sourceEvent, eventData)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish contact creation event"))
		return err
	}

	return nil
}

func publishEventCreateContact(c context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.publishEventCreateContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var err error

	system, err := sourceEvent.ExternalSystem()
	if err != nil {
		return err
	}

	for _, email := range *eventData.ParticipantEmails {
		flowActionEvent := dto.FlowActionEvent{
			ExternalSystemId: system,
			SourceEvent:      sourceEvent,
			Name:             commonenum.ActionCreateContact,
			DataType:         data_fields.ContactCreateEvent{}.Type(),
			Data: data_fields.ContactCreateEvent{
				Email: email,
			},
		}

		pubErr := s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent)
		if pubErr != nil {
			tracing.TraceErr(span, err)
			err = multierr.Append(err, fmt.Errorf("failed to publish contact creation for email %s: %w", email, pubErr))
		}
	}

	return nil
}

func publishCreateMarkdownEvent(c context.Context, s *service.Services, sourceEvent commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.publishCreateMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var mdEvent data_fields.MarkdownEventFields
	var sourceId entity.DataSource

	system, err := sourceEvent.ExternalSystem()
	if err != nil {
		return err
	}

	switch system {
	case enum.Fathom:
		sourceId = entity.DataSourceFathom
	case enum.Grain:
		sourceId = entity.DataSourceGrain
	default:
		return fmt.Errorf("Unuspported source: %v", system)
	}

	mdEvent.Source = &sourceId
	mdEvent.Content = eventData.Content
	mdEvent.CreatedAt = eventData.Timestamp

	flowActionEvent := dto.FlowActionEvent{
		ExternalSystemId: system,
		SourceEvent:      sourceEvent,
		Name:             commonenum.ActionCreateTimelineEvent,
		DataType:         data_fields.MarkdownEventFields{}.Type(),
		Data:             &mdEvent,
	}

	return s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent)
}
