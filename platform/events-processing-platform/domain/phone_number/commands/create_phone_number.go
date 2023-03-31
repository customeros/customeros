package commands

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type CreatePhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *CreatePhoneNumberCommand) error
}

type createPhoneNumberHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewCreatePhoneNumberHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *createPhoneNumberHandler {
	return &createPhoneNumberHandler{log: log, cfg: cfg, es: es}
}

func (c *createPhoneNumberHandler) Handle(ctx context.Context, command *CreatePhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createPhoneNumberHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	}

	if err := phoneNumberAggregate.CreatePhoneNumber(ctx, command.Tenant, command.PhoneNumber); err != nil {
		return err
	}

	//span.LogFields(log.String("order", order.String()))
	return c.es.Save(ctx, phoneNumberAggregate)
}
