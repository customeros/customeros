package handlers

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

func HandleCreateMarkdownEvent(c context.Context, s *service.Services, eventData *data_fields.MarkdownEventFields) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleCreateMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	_, err := s.MarkdownEventService.Save(ctx, nil, nil, *eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		// todo Implement retry logic?
		return err
	}
	return nil
}
