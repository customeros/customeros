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

type UpsertUserCommandHandler interface {
	Handle(ctx context.Context, command *command.UpsertUserCommand) error
}

type upsertUserCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertUserCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertUserCommandHandler {
	return &upsertUserCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertUserCommandHandler) Handle(ctx context.Context, command *command.UpsertUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertUserCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.UserID)
	span.LogFields(log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(userAggregate) {
		command.IsCreateCommand = true
	}
	if err = userAggregate.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, userAggregate)
}
