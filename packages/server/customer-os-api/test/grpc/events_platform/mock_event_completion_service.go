package events_platform

import (
	"context"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockEventCompletionServiceCallbacks struct {
	CompletionEvent func(context.Context, *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error)
}

var eventCompletionServiceCallbacks = &MockEventCompletionServiceCallbacks{}

func SetEventCompletionServiceCallbacks(callbacks *MockEventCompletionServiceCallbacks) {
	eventCompletionServiceCallbacks = callbacks
}

type MockEventCompletionService struct {
	eventcompletionpb.UnimplementedEventCompletionGrpcServiceServer
}

func (MockEventCompletionService) CompletionEvent(context context.Context, proto *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
	if eventCompletionServiceCallbacks.CompletionEvent == nil {
		panic("eventCompletionServiceCallbacks.Completion is not set")
	}
	return eventCompletionServiceCallbacks.CompletionEvent(context, proto)
}
