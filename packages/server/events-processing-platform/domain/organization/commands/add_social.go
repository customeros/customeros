package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type AddSocialCommandHandler interface {
	Handle(ctx context.Context, command *AddSocialCommand) error
}

type addSocialCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewAddSocialCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) AddSocialCommandHandler {
	return &addSocialCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *addSocialCommandHandler) Handle(ctx context.Context, command *AddSocialCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddSocialCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if command.Tenant == "" {
		return eventstore.ErrMissingTenant
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if err = organizationAggregate.AddSocial(ctx, command.Tenant, command.SocialId, command.SocialPlatform, command.SocialUrl, command.Source, command.CreatedAt, command.UpdatedAt); err != nil {
		return err
	}

	return c.es.Save(ctx, organizationAggregate)
}
