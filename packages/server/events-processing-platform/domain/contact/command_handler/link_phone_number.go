package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error
}

type linkPhoneNumberCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkPhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkPhoneNumberCommandHandler {
	return &linkPhoneNumberCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *linkPhoneNumberCommandHandler) Handle(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkPhoneNumberCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	if cmd.Tenant == "" {
		return eventstore.ErrMissingTenant
	}
	if cmd.PhoneNumberId == "" {
		return errors.ErrMissingPhoneNumberId
	}

	contactAggregate, err := aggregate.LoadContactAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = contactAggregate.LinkPhoneNumber(ctx, cmd.Tenant, cmd.PhoneNumberId, cmd.Label, cmd.Primary); err != nil {
		return err
	}
	if cmd.Primary {
		for k, v := range contactAggregate.Contact.PhoneNumbers {
			if k != cmd.PhoneNumberId && v.Primary {
				if err = contactAggregate.SetPhoneNumberNonPrimary(ctx, cmd.Tenant, cmd.PhoneNumberId); err != nil {
					return err
				}
			}
		}
	}

	return h.es.Save(ctx, contactAggregate)
}
