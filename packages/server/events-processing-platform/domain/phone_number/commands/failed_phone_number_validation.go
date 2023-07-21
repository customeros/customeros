package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FailedPhoneNumberValidationCommandHandler interface {
	Handle(ctx context.Context, command *FailedPhoneNumberValidationCommand) error
}

type failedPhoneNumberValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewFailedPhoneNumberValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) FailedPhoneNumberValidationCommandHandler {
	return &failedPhoneNumberValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *failedPhoneNumberValidationCommandHandler) Handle(ctx context.Context, command *FailedPhoneNumberValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "failedPhoneNumberValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	phoneNumberAggregate, err := aggregate.LoadPhoneNumberAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = phoneNumberAggregate.FailedPhoneNumberValidation(ctx, command.Tenant, command.RawPhoneNumber, command.CountryCodeA2, command.ValidationError); err != nil {
		return err
	}
	return h.es.Save(ctx, phoneNumberAggregate)
}
