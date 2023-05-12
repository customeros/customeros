package email_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type PhoneNumberEventHandler struct {
	phoneNumberCommands *commands.PhoneNumberCommands
	log                 logger.Logger
	cfg                 *config.Config
}

type PhoneNumberValidate struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type PhoneNumberValidationResponseV1 struct {
	Error       string `json:"error"`
	PhoneNumber string `json:"phoneNumber"`
}

func (h *PhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	phoneNumberId := aggregate.GetPhoneNumberID(evt.AggregateID, eventData.Tenant)

	// First iteration, just skip validation
	h.phoneNumberCommands.SkipPhoneNumberValidation.Handle(ctx, commands.NewSkippedPhoneNumberValidationCommand(phoneNumberId, eventData.Tenant, "Skipped"))

	return nil
}
