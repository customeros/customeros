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

type RequestRenewalForecastCommandHandler interface {
	Handle(ctx context.Context, command *command.RequestRenewalForecastCommand) error
}

type requestRenewalForecastCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRequestRenewalForecastCommandHandler(log logger.Logger, es eventstore.AggregateStore) RequestRenewalForecastCommandHandler {
	return &requestRenewalForecastCommandHandler{log: log, es: es}
}

func (h *requestRenewalForecastCommandHandler) Handle(ctx context.Context, cmd *command.RequestRenewalForecastCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestRenewalForecastCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		return err
	}

	if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, organizationAggregate)
}
