package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetInteractionEventObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, InteractionEventAggregateType)
}

func LoadInteractionEventAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InteractionEventAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInteractionEventAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	interactionEventAggregate := NewInteractionEventAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, interactionEventAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Bool("AggregateExists", true))
	return interactionEventAggregate, nil
}
