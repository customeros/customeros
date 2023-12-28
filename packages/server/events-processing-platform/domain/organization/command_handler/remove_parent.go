package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
)

type RemoveParentCommandHandler interface {
	Handle(ctx context.Context, command *command.RemoveParentCommand) error
}

type removeParentCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRemoveParentCommandHandler(log logger.Logger, es eventstore.AggregateStore) RemoveParentCommandHandler {
	return &removeParentCommandHandler{log: log, es: es}
}

func (c *removeParentCommandHandler) Handle(ctx context.Context, cmd *command.RemoveParentCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "removeParentCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		return err
	}

	if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, organizationAggregate)
}
