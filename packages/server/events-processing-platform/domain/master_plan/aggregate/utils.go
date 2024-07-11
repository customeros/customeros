package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetMasterPlanObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, MasterPlanAggregateType)
}

func LoadMasterPlanAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*MasterPlanAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadMasterPlanAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	masterPlanAggregate := NewMasterPlanAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, masterPlanAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return masterPlanAggregate, nil
}
