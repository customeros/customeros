package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func OnFlowActionEventCreated(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowActionEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	tenant, event, err := getFlowActionEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	c := model.EventContext{
		Context:      ctx,
		Span:         span,
		Services:     services,
		Tenant:       tenant,
		SourceEvent:  event.SourceEvent,
		SourceSystem: event.ExternalSystemId,
	}

	// determine event handler //
	switch eventData := (*event.Data).(type) {

	case *data_fields.MarkdownEventFields:
		err := handlers.HandleCreateMarkdownEvent(c, eventData)
		if err != nil {
			tracing.TraceErr(c.Span, err)
			return err
		}

	default:
		err := fmt.Errorf("Unsupported flow action event %s", event.Name)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func getFlowActionEvent(ctx context.Context, input any) (string, *dto.FlowActionEvent[any], error) {
	message := input.(*dto.Event)
	tenant := message.Event.Tenant
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return tenant, nil, err
	}

	return tenant, message.Event.Data.(*dto.FlowActionEvent[any]), nil
}
