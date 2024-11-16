package events_platform

import (
	"context"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
)

type MockTenantServiceCallbacks struct {
	AddBillingProfile    func(context.Context, *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error)
	UpdateBillingProfile func(context.Context, *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error)
}

var tenantCallbacks = &MockTenantServiceCallbacks{}

func SetTenantCallbacks(callbacks *MockTenantServiceCallbacks) {
	tenantCallbacks = callbacks
}

type MockTenantService struct {
	tenantpb.UnimplementedTenantGrpcServiceServer
}

func (MockTenantService) AddBillingProfile(context context.Context, proto *tenantpb.AddBillingProfileRequest) (*commonpb.IdResponse, error) {
	if tenantCallbacks.AddBillingProfile == nil {
		panic("tenantCallbacks.AddBillingProfile is not set")
	}
	return tenantCallbacks.AddBillingProfile(context, proto)
}

func (MockTenantService) UpdateBillingProfile(context context.Context, proto *tenantpb.UpdateBillingProfileRequest) (*commonpb.IdResponse, error) {
	if tenantCallbacks.UpdateBillingProfile == nil {
		panic("tenantCallbacks.UpdateBillingProfile is not set")
	}
	return tenantCallbacks.UpdateBillingProfile(context, proto)
}
