package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetContactObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(ContactAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadContactAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*ContactAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContactAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	contactAggregate := NewContactAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, contactAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return nil, err
	}

	if err := eventStore.Load(ctx, contactAggregate); err != nil {
		return nil, err
	}

	return contactAggregate, nil
}
