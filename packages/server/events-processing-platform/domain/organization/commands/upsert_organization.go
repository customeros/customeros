package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *UpsertOrganizationCommand) error
}

type upsertOrganizationCommandHandler struct {
	log    logger.Logger
	cfg    *config.Config
	es     eventstore.AggregateStore
	caches caches.Cache
}

func NewUpsertOrganizationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, caches caches.Cache) UpsertOrganizationCommandHandler {
	return &upsertOrganizationCommandHandler{log: log, cfg: cfg, es: es, caches: caches}
}

func (c *upsertOrganizationCommandHandler) Handle(ctx context.Context, command *UpsertOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertOrganizationCommandHandler.Handle")
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

	orgFields := UpsertOrganizationCommandToOrganizationFields(command)
	orgFields.OrganizationDataFields.Market = data.AdjustOrganizationMarket(orgFields.OrganizationDataFields.Market)
	orgFields.OrganizationDataFields.Industry = adjustIndustryValue(orgFields.OrganizationDataFields.Industry, c.caches)

	if aggregate.IsAggregateNotFound(organizationAggregate) {
		if err = organizationAggregate.CreateOrganization(ctx, orgFields); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		if err = organizationAggregate.UpdateOrganization(ctx, orgFields); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return c.es.Save(ctx, organizationAggregate)
}
