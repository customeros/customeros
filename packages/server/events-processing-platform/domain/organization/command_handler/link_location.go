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
	"github.com/opentracing/opentracing-go/log"
)

type LinkLocationCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkLocationCommand) error
}

type linkLocationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewLinkLocationCommandHandler(log logger.Logger, es eventstore.AggregateStore) LinkLocationCommandHandler {
	return &linkLocationCommandHandler{log: log, es: es}
}

func (c *linkLocationCommandHandler) Handle(ctx context.Context, cmd *command.LinkLocationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkLocationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

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
