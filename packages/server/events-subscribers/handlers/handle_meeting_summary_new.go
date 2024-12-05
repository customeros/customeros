package handlers

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func HandleMeetingSummaryEvent(c context.Context, s *service.Services, eventName commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleMeetingSummaryEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// todo lookup which actions are configured as part of flow for tenant
	// only trigger events for actions that are turned on

	err := publishCreateMarkdownEvent(ctx, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func publishEventCreateContact(c context.Context, s *service.Services, eventName commonenum.FlowEvent, eventData *data_fields.MeetingSummaryEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.publishEventCreateContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var err error

	system, err := eventName.ExternalSystem()
	if err != nil {
		return err
	}

	for _, email := range *eventData.ParticipantEmails {
		flowActionEvent := dto.NewFlowActionEvent(
			commonenum.ActionCreateContact,
			system,
			eventName,
			data_fields.ContactCreateEvent{
				Email: email,
			},
		)

		pubErr := s.RabbitMQService.PublishFlowActionEvent(ctx, flowActionEvent)
		if pubErr != nil {
			tracing.TraceErr(span, err)
			err = multierr.Append(err, fmt.Errorf("failed to publish contact creation for email %s: %w", email, pubErr))
		}
	}

	return nil
}

func publishCreateMarkdownEvent(ctx model.EventContext, event *data_fields.MeetingSummaryEvent) error {
	var mdEvent data_fields.MarkdownEventFields
	var sourceId entity.DataSource

	switch ctx.SourceSystem {
	case enum.Fathom:
		sourceId = entity.DataSourceFathom
	case enum.Grain:
		sourceId = entity.DataSourceGrain
	default:
		return fmt.Errorf("Unuspported source: %v", ctx.SourceSystem)
	}

	mdEvent.Source = &sourceId
	mdEvent.Content = event.Content
	mdEvent.CreatedAt = event.Timestamp

	flowActionEvent := dto.NewFlowActionEvent(
		commonenum.ActionCreateTimelineEvent,
		ctx.SourceSystem,
		ctx.SourceEvent,
		mdEvent,
	)

	return ctx.Services.RabbitMQService.PublishFlowActionEvent(ctx.Context, flowActionEvent)
}
