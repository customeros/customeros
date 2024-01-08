package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadInvoicingCycleAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InvoicingCycleAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInvoicingCycleAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	invoicingCycleAggregate := NewInvoicingCycleAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, invoicingCycleAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return invoicingCycleAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, invoicingCycleAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return invoicingCycleAggregate, nil
}
