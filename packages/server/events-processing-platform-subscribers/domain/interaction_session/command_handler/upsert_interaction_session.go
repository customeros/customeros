package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertInteractionSessionCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertInteractionSessionCommand) error
}

type upsertInteractionSessionCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertInteractionSessionCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertInteractionSessionCommandHandler {
	return &upsertInteractionSessionCommandHandler{log: log, es: es}
}

func (c *upsertInteractionSessionCommandHandler) Handle(ctx context.Context, cmd *command.UpsertInteractionSessionCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertInteractionSessionCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("command", fmt.Sprintf("%+v", cmd)))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		wrappedErr := errors.Wrap(err, "failed validation for UpsertInteractionSessionCommand")
		tracing.TraceErr(span, wrappedErr)
		return wrappedErr
	}

	interactionSessionAggregate, err := aggregate.LoadInteractionSessionAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if interactionSessionAggregate.NotFound() {
		cmd.IsCreateCommand = true
	}
	if err = interactionSessionAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, interactionSessionAggregate)
}
