package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error
}

type upsertPhoneNumberHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertPhoneNumberHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertPhoneNumberCommandHandler {
	return &upsertPhoneNumberHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertPhoneNumberHandler) Handle(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertUserCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	phoneNumberAggregate, err := aggregate.LoadPhoneNumberAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(phoneNumberAggregate) {
		cmd.IsCreateCommand = true
	}
	if err = phoneNumberAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, phoneNumberAggregate)
}
