package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CreateContactCommandHandler interface {
	Handle(ctx context.Context, command *CreateContactCommand) error
}

type createContactHandler struct {
	log        logger.Logger
	cfg        *config.Config
	eventStore eventstore.AggregateStore
}

func NewCreateContactHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *createContactHandler {
	return &createContactHandler{log: log, cfg: cfg, eventStore: es}
}

func (c *createContactHandler) Handle(ctx context.Context, command *CreateContactCommand) error {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "createContactHandler.Handle")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", command.GetAggregateID()))

	//contactAggregate := aggregate.NewContactAggregateWithID(command.AggregateID)
	//err := c.eventStore.Exists(ctx, contactAggregate.GetID())
	//if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
	//	return err
	//}
	//
	//if err := contactAggregate.CreateContact(ctx, command.UUID, command.FirstName, command.LastName); err != nil {
	//	return err
	//}
	//
	////span.LogFields(log.String("contactAggregate", contactAggregate.String()))
	//return c.eventStore.Save(ctx, contactAggregate)
	return nil
}
