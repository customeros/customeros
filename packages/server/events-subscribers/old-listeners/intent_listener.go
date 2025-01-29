package oldlisteners

import (
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

func OnIntentEventCreated(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message, err := validateEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	intentData, err := validateIntentEvent(ctx, message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	ctx = common.SetTenantInContext(ctx, intentData.Tenant)

	return distributeIntentEventToSubscribers(ctx, dependencies, intentData)
}

func distributeIntentEventToSubscribers(ctx context.Context, dep *model.DependencyContainer, intentData *dto.IntentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.distributeIntentEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	switch intentData.IntentType {
	case enum.IntentSupportRequired:
		return dep.CommonServices.SupportAgent.ProcessIntentEvent(ctx, intentData)

	default:
		err := errors.New("IntentType not supported")
		tracing.TraceErr(span, err)
		return err
	}
}

func validateIntentEvent(ctx context.Context, message *dto.Event) (*dto.IntentEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.validateIntentEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	messageData, ok := message.Event.Data.(*dto.IntentEvent)
	if !ok {
		err := errors.New("cannot cast type to dto.IntentEvent")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return messageData, nil
}
