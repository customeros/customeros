package handlers

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

func HandleCreateContact(c context.Context, s *service.Services, eventData *data_fields.ContactCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleCreateContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	_, err := s.ContactService.CreateContactWithOrganizationByEmail(ctx, nil, eventData.Email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
		// todo Implement retry logic?
	}
	return nil
}
