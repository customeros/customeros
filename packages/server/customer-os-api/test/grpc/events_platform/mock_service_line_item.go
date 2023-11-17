package events_platform

import (
	"context"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
)

type MockServiceLineItemServiceCallbacks struct {
	CreateServiceLineItem func(context.Context, *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error)
}

var serviceLineItemCallbacks = &MockServiceLineItemServiceCallbacks{}

func SetServiceLineItemCallbacks(callbacks *MockServiceLineItemServiceCallbacks) {
	serviceLineItemCallbacks = callbacks
}

type MockServiceLineItemService struct {
	servicelineitempb.UnimplementedServiceLineItemGrpcServiceServer
}

func (MockServiceLineItemService) CreateServiceLineItem(context context.Context, proto *servicelineitempb.CreateServiceLineItemGrpcRequest) (*servicelineitempb.ServiceLineItemIdGrpcResponse, error) {
	if serviceLineItemCallbacks.CreateServiceLineItem == nil {
		panic("serviceLineItemCallbacks.CreateServiceLineItem is not set")
	}
	return serviceLineItemCallbacks.CreateServiceLineItem(context, proto)
}
