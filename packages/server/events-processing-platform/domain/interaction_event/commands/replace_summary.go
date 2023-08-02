package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ReplaceSummaryCommandHandler interface {
	Handle(ctx context.Context, command *ReplaceSummaryCommand) error
}

type ReplaceSummaryCommand struct {
	eventstore.BaseCommand
	Summary     string
	ContentType string
	UpdatedAt   *time.Time
}

func NewReplaceSummaryCommand(tenant, interactionEventId, summary, contentType string, updatedAt *time.Time) *ReplaceSummaryCommand {
	return &ReplaceSummaryCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant),
		Summary:     summary,
		ContentType: contentType,
		UpdatedAt:   updatedAt,
	}
}

type replaceSummaryCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewReplaceSummaryCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) ReplaceSummaryCommandHandler {
	return &replaceSummaryCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *replaceSummaryCommandHandler) Handle(ctx context.Context, command *ReplaceSummaryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReplaceSummaryCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.ReplaceSummary(ctx, command.Tenant, command.Summary, command.ContentType, command.UpdatedAt); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, interactionEventAggregate)
}
