package email_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmailEventHandler struct {
	emailCommands *commands.EmailCommands
}

type EmailValidate struct {
	Email string `validate:"required,email"`
}

func (e *EmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	preValidationErr := validator.GetValidator().Struct(EmailValidate{
		Email: eventData.RawEmail,
	})
	if preValidationErr != nil {
		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(evt.GetAggregateID(), eventData.Tenant, preValidationErr.Error()))
	} else {
		// FIXME alexb implement invoking Validatio API
	}

	// FIXME alexb - implement if validation API not reachable add error
	// FIXME alexb - implement if validation API returns error add error
	// FIXME alexb - implement validation API correct validation

	return nil
}
