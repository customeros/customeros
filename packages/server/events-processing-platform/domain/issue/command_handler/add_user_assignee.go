package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type AddUserAssigneeCommandHandler interface {
	Handle(ctx context.Context, cmd *command.AddUserAssigneeCommand) error
}

type addUserAssigneeCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewAddUserAssigneeCommandHandler(log logger.Logger, es eventstore.AggregateStore) AddUserAssigneeCommandHandler {
	return &addUserAssigneeCommandHandler{log: log, es: es}
}

func (c *addUserAssigneeCommandHandler) Handle(ctx context.Context, cmd *command.AddUserAssigneeCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "addUserAssigneeCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("command", fmt.Sprintf("%+v", cmd)))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		wrappedErr := errors.Wrap(err, "failed validation for AddUserAssigneeCommand")
		tracing.TraceErr(span, wrappedErr)
		return wrappedErr
	}

	issueAggregate, err := aggregate.LoadIssueAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = issueAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, issueAggregate)
}
