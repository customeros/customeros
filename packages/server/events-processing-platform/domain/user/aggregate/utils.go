package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetUserObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(UserAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() <= 0
}

func LoadUserAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*UserAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadUserAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	userAggregate := NewUserAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, userAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return userAggregate, nil
		}
	}

	if err := eventStore.Load(ctx, userAggregate); err != nil {
		return nil, err
	}

	return userAggregate, nil
}
