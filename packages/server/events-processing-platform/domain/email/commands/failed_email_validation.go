package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FailEmailValidationCommandHandler interface {
	Handle(ctx context.Context, command *FailedEmailValidationCommand) error
}

type failEmailValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewFailEmailValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *failEmailValidationCommandHandler {
	return &failEmailValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *failEmailValidationCommandHandler) Handle(ctx context.Context, command *FailedEmailValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "failEmailValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, emailAggregate.GetID())

	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	emailAggregate, _ = aggregate.LoadEmailAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err = emailAggregate.FailEmailValidation(ctx, command.Tenant, command.ValidationError); err != nil {
		return err
	}
	return c.es.Save(ctx, emailAggregate)
}
