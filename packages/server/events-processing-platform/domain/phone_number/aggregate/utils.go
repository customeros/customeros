package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetPhoneNumberObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(PhoneNumberAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() <= 0
}

func LoadPhoneNumberAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*PhoneNumberAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadPhoneNumberAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	phoneNumberAggregate := NewPhoneNumberAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return phoneNumberAggregate, nil
		}
	}

	if err := eventStore.Load(ctx, phoneNumberAggregate); err != nil {
		return nil, err
	}

	return phoneNumberAggregate, nil
}
