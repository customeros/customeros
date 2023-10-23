package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpdateRenewalForecastCommandHandler interface {
	Handle(ctx context.Context, command *cmd.UpdateRenewalForecastCommand) error
}

type updateRenewalForecastCommandHandler struct {
	log          logger.Logger
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewUpdateRenewalForecastCommandHandler(log logger.Logger, es eventstore.AggregateStore, repositories *repository.Repositories) UpdateRenewalForecastCommandHandler {
	return &updateRenewalForecastCommandHandler{log: log, es: es, repositories: repositories}
}

func (h *updateRenewalForecastCommandHandler) Handle(ctx context.Context, command *cmd.UpdateRenewalForecastCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateRenewalForecastCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.LoggedInUserId)
	span.LogFields(log.Object("command", command))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "Command details validation failed")
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, command.Tenant, command.ObjectID)
	if err != nil {
		err = errors.Wrap(err, "Organization not found")
		tracing.TraceErr(span, err)
		return err
	}
	if orgDbNode == nil {
		err = errors.New("Organization not found")
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(organizationAggregate) {
		if err = organizationAggregate.CreateOrganization(ctx, &models.OrganizationFields{
			ID:     command.ObjectID,
			Tenant: command.Tenant,
		}, command.LoggedInUserId); err != nil {
			err := errors.Wrap(err, "Error while creating organization")
			tracing.TraceErr(span, err)
			return err
		}
	}

	if err = organizationAggregate.HandleCommand(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, organizationAggregate)
}
