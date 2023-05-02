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

type FailEmailValidationCommandHandler interface {
	Handle(ctx context.Context, command *FailEmailValidationCommand) error
}

type failEmailValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewFailEmailValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *failEmailValidationCommandHandler {
	return &failEmailValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *failEmailValidationCommandHandler) Handle(ctx context.Context, command *FailEmailValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "failEmailValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, emailAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	}

	if err = emailAggregate.FailEmailValidation(ctx, command.Tenant, command.ValidationError); err != nil {
		return err
	}
	return c.es.Save(ctx, emailAggregate)
}
