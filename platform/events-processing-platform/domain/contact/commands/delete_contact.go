package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
)

type DeleteContactCommandHandler interface {
	Handle(ctx context.Context, command *DeleteContactCommand) error
}

type deleteContactCommandHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore eventstore.AggregateStore
}

func NewDeleteContactCommandHandler(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore) *deleteContactCommandHandler {
	return &deleteContactCommandHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (c *deleteContactCommandHandler) Handle(ctx context.Context, command *DeleteContactCommand) error {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "deleteContactCommandHandler.Handle")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	contactAggregate, err := aggregate.LoadContactAggregate(ctx, c.aggregateStore, command.GetAggregateID())
	if err != nil {
		return err
	}

	if err := contactAggregate.DeleteContact(ctx, command.UUID); err != nil {
		return err
	}

	return c.aggregateStore.Save(ctx, contactAggregate)
}
