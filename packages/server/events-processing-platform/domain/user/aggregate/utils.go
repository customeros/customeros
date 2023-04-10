package aggregate

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

// GetUserAggregateID get user aggregate id for eventstoredb
func GetUserAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, string(UserAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadUserAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, aggregateID string) (*UserAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadUserAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", aggregateID))

	userAggregate := NewUserAggregateWithTenantAndID(tenant, aggregateID)

	err := eventStore.Exists(ctx, userAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, err
	}

	if err := eventStore.Load(ctx, userAggregate); err != nil {
		return nil, err
	}

	return userAggregate, nil
}
