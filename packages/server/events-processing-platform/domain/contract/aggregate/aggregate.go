package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	ContractAggregateType eventstore.AggregateType = "contract"
)

type ContractAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Contract *model.Contract
}

func NewContractAggregateWithTenantAndID(tenant, id string) *ContractAggregate {
	contractAggregate := ContractAggregate{}
	contractAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractAggregate.SetWhen(contractAggregate.When)
	contractAggregate.Contract = &model.Contract{}
	contractAggregate.Tenant = tenant

	return &contractAggregate
}

func (a *ContractAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContractCreateV1:
		return a.onContractCreate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *ContractAggregate) onContractCreate(evt eventstore.Event) error {
	var eventData event.ContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.ID = a.ID
	a.Contract.Tenant = a.Tenant
	a.Contract.OrganizationId = eventData.OrganizationId
	a.Contract.Name = eventData.Name
	a.Contract.CreatedByUserId = eventData.CreatedByUserId
	a.Contract.ServiceStartedAt = eventData.ServiceStartedAt
	a.Contract.SignedAt = eventData.SignedAt
	a.Contract.RenewalCycle = eventData.RenewalCycle
	a.Contract.Status = eventData.Status
	a.Contract.CreatedAt = eventData.CreatedAt
	a.Contract.UpdatedAt = eventData.UpdatedAt
	a.Contract.Source = eventData.Source
	if eventData.ExternalSystem.Available() {
		a.Contract.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}

	return nil
}
