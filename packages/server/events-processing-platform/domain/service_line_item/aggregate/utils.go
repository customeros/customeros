package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

// GetServiceLineItemObjectID generates the object ID for a service line item.
func GetServiceLineItemObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getServiceLineItemObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(ServiceLineItemAggregateType)+"-"+tenant+"-", "")
}

// getServiceLineItemObjectUUID generates the UUID for a service line item when the tenant is not known.
func getServiceLineItemObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

// IsAggregateNotFound checks if the provided aggregate is not found.
func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

// LoadServiceLineItemAggregate loads the service line item aggregate from the event store.
func LoadServiceLineItemAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*ServiceLineItemAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadServiceLineItemAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	serviceLineItemAggregate := NewServiceLineItemAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, serviceLineItemAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return serviceLineItemAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, serviceLineItemAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return serviceLineItemAggregate, nil
}
