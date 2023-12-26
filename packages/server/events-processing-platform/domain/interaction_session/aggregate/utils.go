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

func GetInteractionSessionObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(InteractionSessionAggregateType)+"-"+tenant+"-", "")
}

func LoadInteractionSessionAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InteractionSessionAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInteractionSessionAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	interactionEventAggregate := NewInteractionSessionAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, interactionEventAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			span.LogFields(log.Bool("AggregateExists", false))
			return interactionEventAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, interactionEventAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Bool("AggregateExists", true))
	return interactionEventAggregate, nil
}
