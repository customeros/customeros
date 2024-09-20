package mocked_grpc

import (
	"context"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockEventCompletionCallbacks struct {
	NotifyEventProcessed func(ctx context.Context, proto *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error)
}

var eventCompletionServiceCallbacks = &MockEventCompletionCallbacks{}

func SetEventCompletionServiceCallbacks(callbacks *MockEventCompletionCallbacks) {
	eventCompletionServiceCallbacks = callbacks
}

type MockEventCompletionService struct {
	eventcompletionpb.UnimplementedEventCompletionGrpcServiceServer
}

func (MockEventCompletionService) NotifyEventProcessed(context context.Context, proto *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
	if eventCompletionServiceCallbacks.NotifyEventProcessed == nil {
		panic("eventCompletionServiceCallbacks.Completion is not set")
	}
	return eventCompletionServiceCallbacks.NotifyEventProcessed(context, proto)
}
