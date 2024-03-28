package events_platform

import (
	"context"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
)

type MockOfferingServiceCallbacks struct {
	CreateOffering func(context.Context, *offeringpb.CreateOfferingGrpcRequest) (*commonpb.IdResponse, error)
}

var offeringCallbacks = &MockOfferingServiceCallbacks{}

func SetOfferingCallbacks(callbacks *MockOfferingServiceCallbacks) {
	offeringCallbacks = callbacks
}

type MockOfferingService struct {
	offeringpb.UnimplementedOfferingGrpcServiceServer
}

func (MockOfferingService) CreateOffering(context context.Context, proto *offeringpb.CreateOfferingGrpcRequest) (*commonpb.IdResponse, error) {
	if offeringCallbacks.CreateOffering == nil {
		panic("offeringCallbacks.CreateOffering is not set")
	}
	return offeringCallbacks.CreateOffering(context, proto)
}
