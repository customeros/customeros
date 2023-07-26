package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpdateOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *UpdateOrganizationCommand) error
}

type updateOrganizationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpdateOrganizationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpdateOrganizationCommandHandler {
	return &updateOrganizationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *updateOrganizationCommandHandler) Handle(ctx context.Context, command *UpdateOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateOrganizationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	orgFields := &models.OrganizationFields{
		ID:                     command.ObjectID,
		Tenant:                 command.Tenant,
		IgnoreEmptyFields:      command.IgnoreEmptyFields,
		OrganizationDataFields: command.DataFields,
		Source: commonModels.Source{
			SourceOfTruth: command.SourceOfTruth,
		},
		UpdatedAt: command.UpdatedAt,
	}
	if err = organizationAggregate.UpdateOrganization(ctx, orgFields); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, organizationAggregate)
}
