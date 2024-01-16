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

func GetContactObjectID(aggregateID, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(ContactAggregateType)+"-"+tenant+"-", "")
}

func LoadContactAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*ContactAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContactAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	contactAggregate := NewContactAggregateWithTenantAndID(tenant, objectID)

	// ErrAggregateNotFound is an expected error, in which case we return the contractAggregate without any error.
	err := eventStore.Exists(ctx, contactAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return contactAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, contactAggregate); err != nil {
		return nil, err
	}

	return contactAggregate, nil
}
