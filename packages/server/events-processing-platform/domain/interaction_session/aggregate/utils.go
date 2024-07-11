package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetInteractionSessionObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, InteractionSessionAggregateType)
}

func LoadInteractionSessionAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InteractionSessionAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInteractionSessionAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	interactionSessionAggregate := NewInteractionSessionAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, interactionSessionAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return interactionSessionAggregate, nil
}
