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

// GetLocationObjectID get phone number id for eventstoredb
func GetLocationObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(LocationAggregateType)+"-"+tenant+"-", "")
}

func LoadLocationAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*LocationAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadLocationAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	locationAggregate := NewLocationAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, locationAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return locationAggregate, nil
		}
	}

	if err := eventStore.Load(ctx, locationAggregate); err != nil {
		return nil, err
	}

	return locationAggregate, nil
}
