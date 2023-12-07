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

type UpsertOrganizationCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertOrganizationCommand) error
}

type upsertOrganizationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertOrganizationCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertOrganizationCommandHandler {
	return &upsertOrganizationCommandHandler{log: log, es: es}
}

func (c *upsertOrganizationCommandHandler) Handle(ctx context.Context, cmd *command.UpsertOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertOrganizationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		return err
	}

	orgFields := command.UpsertOrganizationCommandToOrganizationFieldsStruct(cmd)

	if aggregate.IsAggregateNotFound(organizationAggregate) {
		cmd.IsCreateCommand = true
		if err = organizationAggregate.CreateOrganization(ctx, orgFields, cmd.LoggedInUserId); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		if err = organizationAggregate.UpdateOrganization(ctx, orgFields, cmd.LoggedInUserId, cmd.FieldsMask); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return c.es.Save(ctx, organizationAggregate)
}
