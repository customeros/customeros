package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertLogEntryCommandHandler interface {
	Handle(ctx context.Context, command *cmd.UpsertLogEntryCommand) error
}

type upsertLogEntryCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertLogEntryCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertLogEntryCommandHandler {
	return &upsertLogEntryCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertLogEntryCommandHandler) Handle(ctx context.Context, command *cmd.UpsertLogEntryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertLogEntryCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.UserID)
	span.LogFields(log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	logEntryAggregate, err := aggregate.LoadLogEntryAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(logEntryAggregate) {
		command.IsCreateCommand = true
	}
	if err = logEntryAggregate.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, logEntryAggregate)
}
