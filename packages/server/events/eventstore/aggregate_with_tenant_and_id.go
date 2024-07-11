package eventstore

import (
	"context"
	"strings"
)

type CommonTenantIdAggregate struct {
	*AggregateBase
	when *func(event Event) error
}

func (a CommonTenantIdAggregate) NotFound() bool {
	return a.GetVersion() < 0
}

func NewCommonAggregateWithTenantAndId(aggregateType AggregateType, tenant, id string) *CommonTenantIdAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewCommonAggregate(aggregateType)
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func (ca *CommonTenantIdAggregate) setWhen(when func(event Event) error) {
	ca.when = &when
}

func GetAggregateWithTenantAndIdObjectID(aggregateID string, aggregateType AggregateType, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+tenant+"-", "")
}

func NewCommonAggregate(aggregateType AggregateType) *CommonTenantIdAggregate {
	commonAggregate := &CommonTenantIdAggregate{}
	base := NewAggregateBase(commonAggregate.When)
	base.SetType(aggregateType)
	commonAggregate.AggregateBase = base
	return commonAggregate
}

func (a *CommonTenantIdAggregate) When(event Event) error {
	if a.when != nil {
		return (*a.when)(event)
	}
	return nil
}

func (a *CommonTenantIdAggregate) SetWhen(when func(event Event) error) {
	a.when = &when
}

func (a *CommonTenantIdAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	return nil, nil
}
