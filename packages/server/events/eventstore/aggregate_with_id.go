package eventstore

import (
	"context"
	"strings"
)

type CommonIdAggregate struct {
	*AggregateBase
	when *func(event Event) error
}

func (a CommonIdAggregate) NotFound() bool {
	return a.GetVersion() < 0
}

func NewCommonAggregateWithId(aggregateType AggregateType, id string) *CommonIdAggregate {
	if id == "" {
		return nil
	}
	aggregate := &CommonIdAggregate{}
	base := NewAggregateBase(aggregate.When)
	base.SetType(aggregateType)
	aggregate.AggregateBase = base
	aggregate.SetID(id)
	return aggregate
}

func (ca *CommonIdAggregate) setWhen(when func(event Event) error) {
	ca.when = &when
}

func GetAggregateWithIdObjectID(aggregateID string, aggregateType AggregateType) string {
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-", "")
}

func LoadCommonAggregateWithId(ctx context.Context, eventStore AggregateStore, aggregateType AggregateType, objectID string) (*CommonIdAggregate, error) {
	aggregate := NewCommonAggregateWithId(aggregateType, objectID)
	err := eventStore.Load(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	return aggregate, nil
}

func (a *CommonIdAggregate) When(event Event) error {
	if a.when != nil {
		return (*a.when)(event)
	}
	return nil
}

func (a *CommonIdAggregate) SetWhen(when func(event Event) error) {
	a.when = &when
}

func (a *CommonIdAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	return nil, nil
}
