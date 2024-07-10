package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RemoveTagCommandHandler interface {
	Handle(ctx context.Context, command *cmd.RemoveTagCommand) error
}

type removeTagCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRemoveTagCommandHandler(log logger.Logger, es eventstore.AggregateStore) RemoveTagCommandHandler {
	return &removeTagCommandHandler{log: log, es: es}
}

func (c *removeTagCommandHandler) Handle(ctx context.Context, command *cmd.RemoveTagCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RemoveTagCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.LoggedInUserId)
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

	if eventstore.IsAggregateNotFound(logEntryAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	} else {
		if err = logEntryAggregate.HandleCommand(ctx, command); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return c.es.Save(ctx, logEntryAggregate)
}
