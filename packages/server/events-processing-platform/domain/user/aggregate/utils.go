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

func GetUserObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getUserObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(UserAggregateType)+"-"+tenant+"-", "")
}

// use this method when tenant is not known
func getUserObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadUserAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*UserAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadUserAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	userAggregate := NewUserAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, userAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return userAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, userAggregate); err != nil {
		return nil, err
	}

	return userAggregate, nil
}
