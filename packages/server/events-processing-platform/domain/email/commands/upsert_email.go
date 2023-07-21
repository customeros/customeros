package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertEmailCommandHandler interface {
	Handle(ctx context.Context, command *UpsertEmailCommand) error
}

type upsertEmailCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertEmailHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertEmailCommandHandler {
	return &upsertEmailCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertEmailCommandHandler) Handle(ctx context.Context, command *UpsertEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertEmailCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if command.Tenant == "" {
		return eventstore.ErrMissingTenant
	}

	emailAggregate, err := aggregate.LoadEmailAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(emailAggregate) {
		if err = emailAggregate.CreateEmail(ctx, command.Tenant, command.RawEmail, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		if err = emailAggregate.UpdateEmail(ctx, command.RawEmail, command.Tenant, command.Source.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	span.LogFields(log.String("Email", emailAggregate.Email.String()))
	return h.es.Save(ctx, emailAggregate)
}
