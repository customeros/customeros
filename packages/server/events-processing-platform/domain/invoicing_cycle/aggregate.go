package invoicing_cycle

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	InvoicingCycleAggregateType eventstore.AggregateType = "invoicing_cycle"
)

type InvoicingCycleAggregate struct {
	*eventstore.CommonTenantIdAggregate
	InvoicingCycle *InvoicingCycle
}

func GetInvoicingCycleObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, InvoicingCycleAggregateType)
}

func LoadInvoicingCycleAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InvoicingCycleAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInvoicingCycleAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	invoicingCycleAggregate := NewInvoicingCycleAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, invoicingCycleAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return invoicingCycleAggregate, nil
}

func NewInvoicingCycleAggregateWithTenantAndID(tenant, id string) *InvoicingCycleAggregate {
	invoicingCycleAggregate := InvoicingCycleAggregate{}
	invoicingCycleAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(InvoicingCycleAggregateType, tenant, id)
	invoicingCycleAggregate.SetWhen(invoicingCycleAggregate.When)
	invoicingCycleAggregate.InvoicingCycle = &InvoicingCycle{}
	invoicingCycleAggregate.Tenant = tenant

	return &invoicingCycleAggregate
}

func (a *InvoicingCycleAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case InvoicingCycleCreateV1:
		return a.onInvoicingCycleCreate(evt)
	case InvoicingCycleUpdateV1:
		return a.onInvoicingCycleUpdate(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *InvoicingCycleAggregate) onInvoicingCycleCreate(evt eventstore.Event) error {
	var eventData InvoicingCycleCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.InvoicingCycle.ID = a.ID
	a.InvoicingCycle.Type = eventData.Type
	a.InvoicingCycle.CreatedAt = eventData.CreatedAt
	a.InvoicingCycle.SourceFields = eventData.SourceFields

	return nil
}

func (a *InvoicingCycleAggregate) onInvoicingCycleUpdate(evt eventstore.Event) error {
	var eventData InvoicingCycleUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.InvoicingCycle.Type = eventData.Type

	return nil
}
