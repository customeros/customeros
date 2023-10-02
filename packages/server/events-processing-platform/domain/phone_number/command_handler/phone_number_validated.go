package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PhoneNumberValidatedCommandHandler interface {
	Handle(ctx context.Context, cmd *command.PhoneNumberValidatedCommand) error
}

type phoneNumberValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewPhoneNumberValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) PhoneNumberValidatedCommandHandler {
	return &phoneNumberValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *phoneNumberValidatedCommandHandler) Handle(ctx context.Context, cmd *command.PhoneNumberValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberValidatedCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.UserID)
	span.LogFields(log.String("ObjectID", cmd.ObjectID))

	phoneNumberAggregate, err := aggregate.LoadPhoneNumberAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = phoneNumberAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return h.es.Save(ctx, phoneNumberAggregate)
}
