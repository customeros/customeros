package events_platform

import (
	"context"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
)

type MockEventStoreServiceCallbacks struct {
	StoreEvent func(context.Context, *eventstorepb.StoreEventGrpcRequest) (*eventstorepb.StoreEventGrpcResponse, error)
}

var eventStoreServiceCallbacks = &MockEventStoreServiceCallbacks{}

func SetEventStoreServiceCallbacks(callbacks *MockEventStoreServiceCallbacks) {
	eventStoreServiceCallbacks = callbacks
}

type MockEventStoreService struct {
	eventstorepb.UnimplementedEventStoreGrpcServiceServer
}

func (MockEventStoreService) StoreEvent(context context.Context, proto *eventstorepb.StoreEventGrpcRequest) (*eventstorepb.StoreEventGrpcResponse, error) {
	if eventStoreServiceCallbacks.StoreEvent == nil {
		panic("eventStoreServiceCallbacks.Store is not set")
	}
	return eventStoreServiceCallbacks.StoreEvent(context, proto)
}
