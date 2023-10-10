package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertContactCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertContactCommand) error
}

type upsertContactCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertContactCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertContactCommandHandler {
	return &upsertContactCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertContactCommandHandler) Handle(ctx context.Context, cmd *command.UpsertContactCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertContactCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.UserID)
	span.LogFields(log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var contactAggregate *aggregate.ContactAggregate
	var err error
	if cmd.IsCreateCommand {
		contactAggregate = aggregate.NewContactAggregateWithTenantAndID(cmd.Tenant, cmd.ObjectID)
	} else {
		contactAggregate, err = aggregate.LoadContactAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if aggregate.IsAggregateNotFound(contactAggregate) {
		cmd.IsCreateCommand = true
	}

	if err = contactAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, contactAggregate)
}
