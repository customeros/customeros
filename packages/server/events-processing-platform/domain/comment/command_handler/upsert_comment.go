package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	span.LogFields(log.String("command", fmt.Sprintf("%+v", cmd)))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed validation for UpsertCommentCommand")
	}

	commentAggregate, err := aggregate.LoadCommentAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
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
