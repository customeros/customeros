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

func GetEmailObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(EmailAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadEmailAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*EmailAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadEmailAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	emailAggregate := NewEmailAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, emailAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return emailAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, emailAggregate); err != nil {
		return nil, err
	}

	return emailAggregate, nil
}
