package events_platform

import (
	"context"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
)

type MockOpportunityServiceCallbacks struct {
	CreateOpportunity        func(context.Context, *opportunitypb.CreateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateOpportunity        func(context.Context, *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateRenewalOpportunity func(context.Context, *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
}

var opportunityCallbacks = &MockOpportunityServiceCallbacks{}

func SetOpportunityCallbacks(callbacks *MockOpportunityServiceCallbacks) {
	opportunityCallbacks = callbacks
}

type MockOpportunityService struct {
	opportunitypb.UnimplementedOpportunityGrpcServiceServer
}

func (MockOpportunityService) CreateOpportunity(context context.Context, proto *opportunitypb.CreateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.CreateOpportunity == nil {
		panic("opportunityCallbacks.CreateOpportunity is not set")
	}
	return opportunityCallbacks.CreateOpportunity(context, proto)
}

func (MockOpportunityService) UpdateOpportunity(context context.Context, proto *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.UpdateOpportunity == nil {
		panic("opportunityCallbacks.UpdateOpportunity is not set")
	}
	return opportunityCallbacks.UpdateOpportunity(context, proto)
}

func (MockOpportunityService) UpdateRenewalOpportunity(context context.Context, proto *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.UpdateRenewalOpportunity == nil {
		panic("opportunityCallbacks.UpdateRenewalOpportunity is not set")
	}
	return opportunityCallbacks.UpdateRenewalOpportunity(context, proto)
}
