package handlers

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func HandleMeetingSummaryEvent(ctx model.EventContext, eventData *data_fields.MeetingSummaryEvent) error {
	ctx.Span, ctx.Context = opentracing.StartSpanFromContext(ctx.Context, "EventHandlers.HandleMeetingSummaryEvent")
	defer ctx.Span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx.Context, ctx.Span)
	tracing.LogObjectAsJson(ctx.Span, "eventData", eventData)

	// todo lookup which actions are configured as part of flow for tenant
	// only trigger events for actions that are turned on

	err := publishCreateMarkdownEvent(ctx, eventData)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return err
	}

	return nil
}

func publishEventCreateContact() error {
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
