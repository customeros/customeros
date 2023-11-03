package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkOrganizationCommand) error
}

type linkOrganizationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewLinkOrganizationCommandHandler(log logger.Logger, es eventstore.AggregateStore) LinkOrganizationCommandHandler {
	return &linkOrganizationCommandHandler{log: log, es: es}
}

func (h *linkOrganizationCommandHandler) Handle(ctx context.Context, cmd *command.LinkOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkOrganizationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	contactAggregate, err := aggregate.LoadContactAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = contactAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, contactAggregate)
}
