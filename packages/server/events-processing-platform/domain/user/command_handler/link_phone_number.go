package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkPhoneNumberCommand) error
}

type linkPhoneNumberCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkPhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkPhoneNumberCommandHandler {
	return &linkPhoneNumberCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkPhoneNumberCommandHandler) Handle(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkPhoneNumberCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = userAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, userAggregate)
}
