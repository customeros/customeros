package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
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

type UpdateRenewalLikelihoodCommandHandler interface {
	Handle(ctx context.Context, command *command.UpdateRenewalLikelihoodCommand) error
}

type updateRenewalLikelihoodCommandHandler struct {
	log          logger.Logger
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewUpdateRenewalLikelihoodCommandHandler(log logger.Logger, es eventstore.AggregateStore, repositories *repository.Repositories) UpdateRenewalLikelihoodCommandHandler {
	return &updateRenewalLikelihoodCommandHandler{log: log, es: es, repositories: repositories}
}

func (h *updateRenewalLikelihoodCommandHandler) Handle(ctx context.Context, cmd *command.UpdateRenewalLikelihoodCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateRenewalLikelihoodCommandHandler.Handle")
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

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, cmd.Tenant, cmd.ObjectID)
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
			ID:     cmd.ObjectID,
			Tenant: cmd.Tenant,
		}, cmd.LoggedInUserId); err != nil {
			err := errors.Wrap(err, "Error while creating organization")
			tracing.TraceErr(span, err)
			return err
		}
	}

	if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, organizationAggregate)
}
