package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetOpportunityObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, OpportunityAggregateType)
}

func LoadOpportunityAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*OpportunityAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOpportunityAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	opportunityAggregate := NewOpportunityAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, opportunityAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return opportunityAggregate, nil
}
