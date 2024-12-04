package handlers

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func HandleCreateMarkdownEvent(ctx model.EventContext, eventData *data_fields.MarkdownEventFields) error {
	ctx.Span, ctx.Context = opentracing.StartSpanFromContext(ctx.Context, "EventHandlers.HandleMeetingSummaryEvent")
	defer ctx.Span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx.Context, ctx.Span)
	tracing.LogObjectAsJson(ctx.Span, "eventData", eventData)

	_, err := ctx.Services.MarkdownEventService.Save(ctx.Context, nil, nil, *eventData)
	if err != nil {
		tracing.TraceErr(ctx.Span, err)
		return err
	}

	return nil
}
