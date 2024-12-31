package events_platform

import (
	"context"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
)

type MockServiceLineItemServiceCallbacks struct {
	DeleteServiceLineItem func(context.Context, *servicelineitempb.DeleteServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error)
	CloseServiceLineItem  func(context.Context, *servicelineitempb.CloseServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error)
}

var serviceLineItemCallbacks = &MockServiceLineItemServiceCallbacks{}

func SetServiceLineItemCallbacks(callbacks *MockServiceLineItemServiceCallbacks) {
	serviceLineItemCallbacks = callbacks
}

type MockServiceLineItemService struct {
	servicelineitempb.UnimplementedServiceLineItemGrpcServiceServer
}

func (MockServiceLineItemService) DeleteServiceLineItem(context context.Context, proto *servicelineitempb.DeleteServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	if serviceLineItemCallbacks.DeleteServiceLineItem == nil {
		panic("serviceLineItemCallbacks.DeleteServiceLineItem is not set")
	}
	return serviceLineItemCallbacks.DeleteServiceLineItem(context, proto)
}

func (MockServiceLineItemService) CloseServiceLineItem(context context.Context, proto *servicelineitempb.CloseServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	if serviceLineItemCallbacks.CloseServiceLineItem == nil {
		panic("serviceLineItemCallbacks.CloseServiceLineItem is not set")
	}
	return serviceLineItemCallbacks.CloseServiceLineItem(context, proto)
}
