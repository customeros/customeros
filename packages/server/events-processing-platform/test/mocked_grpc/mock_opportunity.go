package mocked_grpc

import (
	"context"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
)

type MockOpportunityServiceCallbacks struct {
	UpdateRenewalOpportunity              func(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
	UpdateRenewalOpportunityNextCycleDate func(ctx context.Context, proto *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error)
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
