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

type LinkEmailCommandHandler interface {
	Handle(ctx context.Context, cmd *command.LinkEmailCommand) error
}

type linkEmailCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkEmailCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkEmailCommandHandler {
	return &linkEmailCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *linkEmailCommandHandler) Handle(ctx context.Context, cmd *command.LinkEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LinkEmailCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	if cmd.Tenant == "" {
		return eventstore.ErrMissingTenant
	}
	if cmd.EmailId == "" {
		return errors.ErrMissingEmailId
	}

	contactAggregate, err := aggregate.LoadContactAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = contactAggregate.LinkEmail(ctx, cmd.Tenant, cmd.EmailId, cmd.Label, cmd.Primary); err != nil {
		return err
	}
	if cmd.Primary {
		for k, v := range contactAggregate.Contact.Emails {
			if k != cmd.EmailId && v.Primary {
				if err = contactAggregate.SetEmailNonPrimary(ctx, cmd.Tenant, cmd.EmailId); err != nil {
					return err
				}
			}
		}
	}

	span.LogFields(log.String("Contact", contactAggregate.Contact.String()))
	return h.es.Save(ctx, contactAggregate)
}
