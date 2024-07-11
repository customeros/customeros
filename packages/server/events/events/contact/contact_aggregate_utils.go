package contact

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetContactObjectID(aggregateID, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, ContactAggregateType)
}

func LoadContactAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, options eventstore.LoadAggregateOptions) (*ContactAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContactAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	contactAggregate := NewContactAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, contactAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return contactAggregate, nil
}
