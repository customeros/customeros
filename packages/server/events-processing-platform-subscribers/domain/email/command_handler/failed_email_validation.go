package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FailEmailValidationCommandHandler interface {
	Handle(ctx context.Context, cmd *command.FailedEmailValidationCommand) error
}

type failEmailValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewFailEmailValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) FailEmailValidationCommandHandler {
	return &failEmailValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *failEmailValidationCommandHandler) Handle(ctx context.Context, cmd *command.FailedEmailValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "failEmailValidationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("ObjectID", cmd.ObjectID))

	emailAggregate, err := aggregate.LoadEmailAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = emailAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, emailAggregate)
}
