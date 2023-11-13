package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	ServiceLineItemAggregateType eventstore.AggregateType = "service_line_item"
)

type ServiceLineItemAggregate struct {
	*aggregate.CommonTenantIdAggregate
	ServiceLineItem *model.ServiceLineItem
}

func NewServiceLineItemAggregateWithTenantAndID(tenant, id string) *ServiceLineItemAggregate {
	serviceLineItemAggregate := ServiceLineItemAggregate{}
	serviceLineItemAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ServiceLineItemAggregateType, tenant, id)
	serviceLineItemAggregate.SetWhen(serviceLineItemAggregate.When)
	serviceLineItemAggregate.ServiceLineItem = &model.ServiceLineItem{}
	serviceLineItemAggregate.Tenant = tenant

	return &serviceLineItemAggregate
}

func (a *ServiceLineItemAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ServiceLineItemCreateV1:
		return a.onServiceLineItemCreate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *ServiceLineItemAggregate) onServiceLineItemCreate(evt eventstore.Event) error {
	var eventData event.ServiceLineItemCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.ServiceLineItem.ID = a.ID
	a.ServiceLineItem.ContractId = eventData.ContractId
	a.ServiceLineItem.Billed = eventData.Billed
	a.ServiceLineItem.Licenses = eventData.Licenses
	a.ServiceLineItem.Price = eventData.Price
	a.ServiceLineItem.Description = eventData.Description
	a.ServiceLineItem.CreatedAt = eventData.CreatedAt
	a.ServiceLineItem.UpdatedAt = eventData.UpdatedAt
	a.ServiceLineItem.Source = eventData.Source

	return nil
}
