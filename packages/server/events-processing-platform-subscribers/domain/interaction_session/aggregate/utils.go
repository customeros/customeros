package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetInteractionSessionObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, InteractionSessionAggregateType)
}

func LoadInteractionSessionAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InteractionSessionAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInteractionSessionAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	interactionSessionAggregate := NewInteractionSessionAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, interactionSessionAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return interactionSessionAggregate, nil
}
