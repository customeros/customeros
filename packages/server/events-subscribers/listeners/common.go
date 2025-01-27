package listeners

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

func validateEvent(ctx context.Context, input any) (*dto.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.ValidateEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := input.(*dto.Event)
	if !ok {
		err := errors.New("unable to cast to event type")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
