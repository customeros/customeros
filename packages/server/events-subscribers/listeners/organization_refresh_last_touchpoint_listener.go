package listeners

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

func OnRequestLastTouchpointRefresh(ctx context.Context, services *service.CommonServices, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnRequestLastTouchpointRefresh")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	organizationId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	err := services.OrganizationService.RefreshLastTouchpoint(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}
