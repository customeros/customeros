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
	"github.com/pkg/errors"
)

type ShowOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *cmd.ShowOrganizationCommand) error
}

type showOrganizationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewShowOrganizationCommandHandler(log logger.Logger, es eventstore.AggregateStore) ShowOrganizationCommandHandler {
	return &showOrganizationCommandHandler{log: log, es: es}
}

func (h *showOrganizationCommandHandler) Handle(ctx context.Context, cmd *cmd.ShowOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ShowOrganizationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID, eventstore.NewLoadAggregateOptions())
	if err != nil {
		return err
	}

	if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "Handle command failed")
	}

	return h.es.Save(ctx, organizationAggregate)
}
