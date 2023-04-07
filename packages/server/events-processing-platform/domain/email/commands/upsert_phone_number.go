package commands

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertEmailCommandHandler interface {
	Handle(ctx context.Context, command *UpsertEmailCommand) error
}

type upsertEmailHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertEmailHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertEmailHandler {
	return &upsertEmailHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertEmailHandler) Handle(ctx context.Context, command *UpsertEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertEmailHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, emailAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	} else if err != nil && errors.Is(err, esdb.ErrStreamNotFound) {
		if err = emailAggregate.CreateEmail(ctx, command.Tenant, command.RawEmail, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		emailAggregate, _ = aggregate.LoadEmailAggregate(ctx, c.es, command.Tenant, command.AggregateID)
		if err = emailAggregate.UpdateEmail(ctx, command.Tenant, command.Source.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	span.LogFields(log.String("Email", emailAggregate.Email.String()))
	return c.es.Save(ctx, emailAggregate)
}
