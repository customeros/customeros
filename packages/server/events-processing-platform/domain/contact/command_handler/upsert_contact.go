package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertContactCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertContactCommand) error
}

type upsertContactCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertContactCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertContactCommandHandler {
	return &upsertContactCommandHandler{log: log, es: es}
}

func (h *upsertContactCommandHandler) Handle(ctx context.Context, cmd *command.UpsertContactCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertContactCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	var contactAggregate *aggregate.ContactAggregate
	var err error
	if cmd.IsCreateCommand {
		contactAggregate = aggregate.NewContactAggregateWithTenantAndID(cmd.Tenant, cmd.ObjectID)
	} else {
		contactAggregate, err = aggregate.LoadContactAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID, *eventstore.NewLoadAggregateOptions())
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if eventstore.IsAggregateNotFound(contactAggregate) {
		cmd.IsCreateCommand = true
	}

	if err = contactAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, contactAggregate)
}
