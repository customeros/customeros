package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertInteractionEventCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertInteractionEventCommand) error
}

type upsertInteractionEventCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertInteractionEventCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertInteractionEventCommandHandler {
	return &upsertInteractionEventCommandHandler{log: log, es: es}
}

func (c *upsertInteractionEventCommandHandler) Handle(ctx context.Context, cmd *command.UpsertInteractionEventCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertInteractionEventCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("command", fmt.Sprintf("%+v", cmd)))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		wrappedErr := errors.Wrap(err, "failed validation for UpsertInteractionEventCommand")
		tracing.TraceErr(span, wrappedErr)
		return wrappedErr
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if interactionEventAggregate.NotFound() {
		cmd.IsCreateCommand = true
	}
	if err = interactionEventAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, interactionEventAggregate)
}
