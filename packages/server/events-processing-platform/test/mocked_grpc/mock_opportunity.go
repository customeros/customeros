package mocked_grpc

import (
	"context"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
)

type MockOpportunityServiceCallbacks struct {
	CreateRenewalOpportunity              func(ctx context.Context, proto *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateRenewalOpportunity              func(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateRenewalOpportunityNextCycleDate func(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateOpportunity                     func(ctx context.Context, proto *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
}

var opportunityCallbacks = &MockOpportunityServiceCallbacks{}

func SetOpportunityCallbacks(callbacks *MockOpportunityServiceCallbacks) {
	opportunityCallbacks = callbacks
}

type MockOpportunityService struct {
	opportunitypb.UnimplementedOpportunityGrpcServiceServer
}

func (MockOpportunityService) UpdateRenewalOpportunity(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.UpdateRenewalOpportunity == nil {
		panic("opportunityCallbacks.UpdateRenewalOpportunity is not set")
	}
	return opportunityCallbacks.UpdateRenewalOpportunity(ctx, proto)
}

func (MockOpportunityService) UpdateRenewalOpportunityNextCycleDate(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.UpdateRenewalOpportunityNextCycleDate == nil {
		panic("opportunityCallbacks.UpdateRenewalOpportunityNextCycleDate is not set")
	}
	return opportunityCallbacks.UpdateRenewalOpportunityNextCycleDate(ctx, proto)
}

func (MockOpportunityService) CreateRenewalOpportunity(ctx context.Context, proto *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.CreateRenewalOpportunity == nil {
		panic("opportunityCallbacks.CreateRenewalOpportunity is not set")
	}
	return opportunityCallbacks.CreateRenewalOpportunity(ctx, proto)
}

func (MockOpportunityService) UpdateOpportunity(ctx context.Context, proto *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	if opportunityCallbacks.UpdateOpportunity == nil {
		panic("opportunityCallbacks.UpdateOpportunity is not set")
	}
	return opportunityCallbacks.UpdateOpportunity(ctx, proto)
}
