package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type RequestNextCycleDateCommandHandler interface {
	Handle(ctx context.Context, command *cmd.RequestNextCycleDateCommand) error
}

type requestNextCycleDateCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRequestNextCycleDateCommandHandler(log logger.Logger, es eventstore.AggregateStore) RequestNextCycleDateCommandHandler {
	return &requestNextCycleDateCommandHandler{log: log, es: es}
}

func (h *requestNextCycleDateCommandHandler) Handle(ctx context.Context, command *cmd.RequestNextCycleDateCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestNextCycleDateCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.LoggedInUserId)
	span.LogFields(log.Object("command", command))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = organizationAggregate.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "Handle command failed")
	}

	return h.es.Save(ctx, organizationAggregate)
}
