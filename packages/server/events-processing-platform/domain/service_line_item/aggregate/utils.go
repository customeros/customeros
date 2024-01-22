package aggregate

import (
	"context"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// GetServiceLineItemObjectID generates the object ID for a service line item.
func GetServiceLineItemObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, ServiceLineItemAggregateType)
}

// getServiceLineItemObjectUUID generates the UUID for a service line item when the tenant is not known.
func getServiceLineItemObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

// LoadServiceLineItemAggregate loads the service line item aggregate from the event store.
func LoadServiceLineItemAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*ServiceLineItemAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadServiceLineItemAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	serviceLineItemAggregate := NewServiceLineItemAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, serviceLineItemAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return serviceLineItemAggregate, nil
}
