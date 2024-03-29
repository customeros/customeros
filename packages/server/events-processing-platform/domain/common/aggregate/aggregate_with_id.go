package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"strings"
)

type CommonIdAggregate struct {
	*eventstore.AggregateBase
	when *func(event eventstore.Event) error
}

func (a CommonIdAggregate) NotFound() bool {
	return a.GetVersion() < 0
}

func NewCommonAggregateWithId(aggregateType eventstore.AggregateType, id string) *CommonIdAggregate {
	if id == "" {
		return nil
	}
	aggregate := &CommonIdAggregate{}
	base := eventstore.NewAggregateBase(aggregate.When)
	base.SetType(aggregateType)
	aggregate.AggregateBase = base
	aggregate.SetID(id)
	return aggregate
}

func (ca *CommonIdAggregate) setWhen(when func(event eventstore.Event) error) {
	ca.when = &when
}

func GetAggregateWithIdObjectID(aggregateID string, aggregateType eventstore.AggregateType) string {
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-", "")
}

func LoadCommonAggregateWithId(ctx context.Context, eventStore eventstore.AggregateStore, aggregateType eventstore.AggregateType, objectID string) (*CommonIdAggregate, error) {
	aggregate := NewCommonAggregateWithId(aggregateType, objectID)
	err := eventStore.Load(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	return aggregate, nil
}

func (a *CommonIdAggregate) When(event eventstore.Event) error {
	if a.when != nil {
		return (*a.when)(event)
	}
	return nil
}

func (a *CommonIdAggregate) SetWhen(when func(event eventstore.Event) error) {
	a.when = &when
}

func (a *CommonIdAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	return nil, nil
}
