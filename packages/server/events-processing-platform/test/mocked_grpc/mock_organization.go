package mocked_grpc

import (
	"context"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
)

type MockOrganizationServiceCallbacks struct {
	RefreshLastTouchpoint func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RefreshRenewalSummary func(ctx context.Context, proto *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RefreshArr            func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
}

var organizationCallbacks = &MockOrganizationServiceCallbacks{}

func SetOrganizationCallbacks(callbacks *MockOrganizationServiceCallbacks) {
	organizationCallbacks = callbacks
}

type MockOrganizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
}

func (MockOrganizationService) RefreshLastTouchpoint(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshLastTouchpoint == nil {
		panic("organizationCallbacks.RefreshLastTouchpoint is not set")
	}
	return organizationCallbacks.RefreshLastTouchpoint(ctx, proto)
}

func (MockOrganizationService) RefreshRenewalSummary(ctx context.Context, proto *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshRenewalSummary == nil {
		panic("organizationCallbacks.RefreshRenewalSummary is not set")
	}
	return organizationCallbacks.RefreshRenewalSummary(ctx, proto)
}

func (MockOrganizationService) RefreshArr(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshArr == nil {
		panic("organizationCallbacks.RefreshArr is not set")
	}
	return organizationCallbacks.RefreshArr(ctx, proto)
}
