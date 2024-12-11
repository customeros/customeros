package handlers

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

func HandleCreateMarkdownEvent(c context.Context, s *service.Services, eventData *data_fields.MarkdownEventFields, flowExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleCreateMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// create markdown event on timeline
	_, err := s.MarkdownEventService.Save(ctx, nil, nil, *eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		// todo Implement retry logic?
		return err
	}

	// create any contacts that don't exist
	_, err = s.ContactService.CreateContactWithOrganizationByEmail(ctx, nil, eventData.Email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
		// todo Implement retry logic?
	}

	// write action execution to db

	// fire action completed event

	// if failure anywhere, fire action failed event

	return nil
}
