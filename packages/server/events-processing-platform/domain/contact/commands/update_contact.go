package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type UpdateContactCommandHandler interface {
	Handle(ctx context.Context, command *UpdateContactCommand) error
}

type updateContactCmdHandler struct {
	log            logger.Logger
	cfg            *config.Config
	aggregateStore eventstore.AggregateStore
}

func NewUpdateContactCmdHandler(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore) UpdateContactCommandHandler {
	return &updateContactCmdHandler{log: log, cfg: cfg, aggregateStore: aggregateStore}
}

func (cmdHandler *updateContactCmdHandler) Handle(ctx context.Context, command *UpdateContactCommand) error {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "updateContactCmdHandler.Handle")
	//defer span.Finish()
	//span.LogFields(log.String("ObjectID", command.GetAggregateID()))

	//contact, err := aggregate.LoadContactAggregate(ctx, cmdHandler.aggregateStore, command.GetAggregateID())
	//if err != nil {
	//	return err
	//}
	//
	//if err := contact.UpdateContact(ctx, command.UUID, command.FirstName, command.LastName); err != nil {
	//	return err
	//}
	//
	//return cmdHandler.aggregateStore.Save(ctx, contact)
	return nil
}
