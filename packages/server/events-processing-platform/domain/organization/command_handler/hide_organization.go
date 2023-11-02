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

type HideOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *cmd.HideOrganizationCommand) error
}

type hideOrganizationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewHideOrganizationCommandHandler(log logger.Logger, es eventstore.AggregateStore) HideOrganizationCommandHandler {
	return &hideOrganizationCommandHandler{log: log, es: es}
}

func (h *hideOrganizationCommandHandler) Handle(ctx context.Context, cmd *cmd.HideOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HideOrganizationCommandHandler.Handle")
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
		return errors.Wrap(err, "Handle command failed")
	}

	return h.es.Save(ctx, organizationAggregate)
}
