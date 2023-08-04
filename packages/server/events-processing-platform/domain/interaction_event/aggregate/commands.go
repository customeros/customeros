package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *InteractionEventAggregate) RequestSummary(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.RequestSummary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewInteractionEventRequestSummaryEvent(a, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventRequestSummaryEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) ReplaceSummary(ctx context.Context, tenant, summary, contentType string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.ReplaceSummary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())

	event, err := events.NewInteractionEventReplaceSummaryEvent(a, tenant, summary, contentType, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventReplaceSummaryEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) RequestActionItems(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.RequestActionItems")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewInteractionEventRequestActionItemsEvent(a, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventRequestActionItemsEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) ReplaceActionItems(ctx context.Context, tenant string, actionItems []string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.ReplaceActionItems")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())

	event, err := events.NewInteractionEventReplaceActionItemsEvent(a, tenant, actionItems, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventReplaceActionItemsEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
