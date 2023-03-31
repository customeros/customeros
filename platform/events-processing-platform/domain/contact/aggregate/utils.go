package aggregate

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

//func GetContactAggregateID(eventAggregateID string) string {
//	return strings.ReplaceAll(eventAggregateID, "order-", "")
//}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadContactAggregate(ctx context.Context, eventStore eventstore.AggregateStore, aggregateID string) (*ContactAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContactAggregate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", aggregateID))

	contactAggregate := NewContactAggregateWithID(aggregateID)

	err := eventStore.Exists(ctx, contactAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, err
	}

	if err := eventStore.Load(ctx, contactAggregate); err != nil {
		return nil, err
	}

	return contactAggregate, nil
}
