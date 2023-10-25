package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"strings"
)

type CommonTenantIdAggregate struct {
	*eventstore.AggregateBase
	when *func(event eventstore.Event) error
}

func (a CommonTenantIdAggregate) NotFound() bool {
	return a.GetVersion() < 0
}

func NewCommonAggregateWithTenantAndId(aggregateType eventstore.AggregateType, tenant, id string) *CommonTenantIdAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewCommonAggregate(aggregateType)
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func (ca *CommonTenantIdAggregate) setWhen(when func(event eventstore.Event) error) {
	ca.when = &when
}

func GetAggregateWithTenantAndIdObjectID(aggregateID string, aggregateType eventstore.AggregateType, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+tenant+"-", "")
}

func LoadCommonAggregateWithTenantAndId(ctx context.Context, eventStore eventstore.AggregateStore, aggregateType eventstore.AggregateType, tenant, objectID string) (*CommonTenantIdAggregate, error) {
	aggregate := NewCommonAggregateWithTenantAndId(aggregateType, tenant, objectID)
	err := eventStore.Load(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	return aggregate, nil
}

func NewCommonAggregate(aggregateType eventstore.AggregateType) *CommonTenantIdAggregate {
	commonAggregate := &CommonTenantIdAggregate{}
	base := eventstore.NewAggregateBase(commonAggregate.When)
	base.SetType(aggregateType)
	commonAggregate.AggregateBase = base
	return commonAggregate
}

func (a *CommonTenantIdAggregate) When(event eventstore.Event) error {
	if a.when != nil {
		return (*a.when)(event)
	}
	return nil
}

func (a *CommonTenantIdAggregate) SetWhen(when func(event eventstore.Event) error) {
	a.when = &when
}
