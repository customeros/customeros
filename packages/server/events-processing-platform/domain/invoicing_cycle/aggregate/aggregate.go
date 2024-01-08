package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	InvoicingCycleAggregateType eventstore.AggregateType = "invoicing_cycle"
)

type InvoicingCycleAggregate struct {
	*aggregate.CommonTenantIdAggregate
	InvoicingCycle *model.InvoicingCycle
}

func NewInvoicingCycleAggregateWithTenantAndID(tenant, id string) *InvoicingCycleAggregate {
	invoicingCycleAggregate := InvoicingCycleAggregate{}
	invoicingCycleAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(InvoicingCycleAggregateType, tenant, id)
	invoicingCycleAggregate.SetWhen(invoicingCycleAggregate.When)
	invoicingCycleAggregate.InvoicingCycle = &model.InvoicingCycle{}
	invoicingCycleAggregate.Tenant = tenant

	return &invoicingCycleAggregate
}

func (a *InvoicingCycleAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.InvoicingCycleCreateV1:
		return a.onInvoicingCycleCreate(evt)
	case event.InvoicingCycleUpdateV1:
		return a.onInvoicingCycleUpdate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *InvoicingCycleAggregate) onInvoicingCycleCreate(evt eventstore.Event) error {
	var eventData event.InvoicingCycleCreateEvent
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
	var eventData event.InvoicingCycleUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.InvoicingCycle.Type = eventData.Type

	return nil
}
