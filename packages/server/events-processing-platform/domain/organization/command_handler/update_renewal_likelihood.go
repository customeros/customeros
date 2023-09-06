package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
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
	cfg          *config.Config
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewUpdateRenewalLikelihoodCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, repositories *repository.Repositories) UpdateRenewalLikelihoodCommandHandler {
	return &updateRenewalLikelihoodCommandHandler{log: log, cfg: cfg, es: es, repositories: repositories}
}

func (h *updateRenewalLikelihoodCommandHandler) Handle(ctx context.Context, command *command.UpdateRenewalLikelihoodCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateRenewalLikelihoodCommandHandler.Handle")
	defer span.Finish()
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
		}, command.UserID); err != nil {
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
