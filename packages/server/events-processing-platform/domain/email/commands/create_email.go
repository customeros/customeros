package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateEmailCommandHandler interface {
	Handle(ctx context.Context, command *CreateEmailCommand) error
}

type createEmailCommandHandler struct {
	log          logger.Logger
	cfg          *config.Config
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewCreateEmailCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) CreateEmailCommandHandler {
	return &createEmailCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *createEmailCommandHandler) Handle(ctx context.Context, command *CreateEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateEmailCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	emailAggregate, err := aggregate.LoadEmailAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = emailAggregate.CreateEmail(ctx, command.Tenant, command.Email, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
		return err
	}

	return h.es.Save(ctx, emailAggregate)
}
