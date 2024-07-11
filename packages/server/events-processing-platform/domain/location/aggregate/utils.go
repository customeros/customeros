package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetLocationObjectID get location id for eventstoredb
func GetLocationObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, LocationAggregateType)
}

func LoadLocationAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*LocationAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadLocationAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	locationAggregate := NewLocationAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, locationAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return locationAggregate, nil
}
