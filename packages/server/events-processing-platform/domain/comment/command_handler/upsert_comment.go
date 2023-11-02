package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertCommentCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertCommentCommand) error
}

type upsertCommentCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertCommentCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertCommentCommandHandler {
	return &upsertCommentCommandHandler{log: log, es: es}
}

func (c *upsertCommentCommandHandler) Handle(ctx context.Context, cmd *command.UpsertCommentCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertCommentCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	commentAggregate, err := aggregate.LoadCommentAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		return err
	}

	if aggregate.IsAggregateNotFound(commentAggregate) {
		cmd.IsCreateCommand = true
	}
	if err = commentAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, commentAggregate)
}
