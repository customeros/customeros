package eventstore

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type testAggregateStore struct {
	aggregateMap map[string][]eventstore.Event
}

func NewTestAggregateStore() *testAggregateStore {
	return &testAggregateStore{aggregateMap: make(map[string][]eventstore.Event)}
}

func (as *testAggregateStore) Load(ctx context.Context, aggregate eventstore.Aggregate) error {
	if _, ok := as.aggregateMap[aggregate.GetID()]; !ok {
		return eventstore.ErrAggregateNotFound
	}
	for _, event := range as.aggregateMap[aggregate.GetID()] {
		if err := aggregate.RaiseEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (as *testAggregateStore) Save(ctx context.Context, aggregate eventstore.Aggregate) error {
	if _, ok := as.aggregateMap[aggregate.GetID()]; !ok {
		as.aggregateMap[aggregate.GetID()] = make([]eventstore.Event, 0)
	}

	for _, event := range aggregate.GetUncommittedEvents() {
		as.aggregateMap[aggregate.GetID()] = append(as.aggregateMap[aggregate.GetID()], event)
	}
	version := 0
	for i := 0; i < len(as.aggregateMap[aggregate.GetID()]); i++ {
		as.aggregateMap[aggregate.GetID()][i].Version = int64(version)
		version++
	}

	return nil
}

func (as *testAggregateStore) Exists(ctx context.Context, aggregateID string) error {
	if _, ok := as.aggregateMap[aggregateID]; !ok {
		return eventstore.ErrAggregateNotFound
	}
	return nil
}

func (as *testAggregateStore) GetEventMap() map[string][]eventstore.Event {
	return as.aggregateMap
}
