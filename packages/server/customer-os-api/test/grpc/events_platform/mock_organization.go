package events_platform

import (
	"context"

	organizationpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
)

type MockOrganizationServiceCallbacks struct {
	RefreshArr                       func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RefreshRenewalSummary            func(ctx context.Context, proto *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	UpdateOrganizationOwner          func(ctx context.Context, proto *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	CreateBillingProfile             func(ctx context.Context, proto *organizationpb.CreateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
	UpdateBillingProfile             func(ctx context.Context, proto *organizationpb.UpdateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
	LinkEmailToBillingProfile        func(ctx context.Context, proto *organizationpb.LinkEmailToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
	UnlinkEmailFromBillingProfile    func(ctx context.Context, proto *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
	LinkLocationToBillingProfile     func(ctx context.Context, proto *organizationpb.LinkLocationToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
	UnlinkLocationFromBillingProfile func(ctx context.Context, proto *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
}

var organizationCallbacks = &MockOrganizationServiceCallbacks{}

func SetOrganizationCallbacks(callbacks *MockOrganizationServiceCallbacks) {
	organizationCallbacks = callbacks
}

type MockOrganizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
}

func (MockOrganizationService) UpdateOrganizationOwner(context context.Context, proto *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.UpdateOrganizationOwner == nil {
		panic("organizationCallbacks.UpdateOrganizationOwner is not set")
	}
	return organizationCallbacks.UpdateOrganizationOwner(context, proto)
}

func (MockOrganizationService) CreateBillingProfile(context context.Context, proto *organizationpb.CreateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.CreateBillingProfile == nil {
		panic("organizationCallbacks.CreateBillingProfile is not set")
	}
	return organizationCallbacks.CreateBillingProfile(context, proto)
}

func (MockOrganizationService) UpdateBillingProfile(context context.Context, proto *organizationpb.UpdateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.UpdateBillingProfile == nil {
		panic("organizationCallbacks.UpdateBillingProfile is not set")
	}
	return organizationCallbacks.UpdateBillingProfile(context, proto)
}

func (MockOrganizationService) LinkEmailToBillingProfile(context context.Context, proto *organizationpb.LinkEmailToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.LinkEmailToBillingProfile == nil {
		panic("organizationCallbacks.LinkEmailToBillingProfile is not set")
	}
	return organizationCallbacks.LinkEmailToBillingProfile(context, proto)
}

func (MockOrganizationService) UnlinkEmailFromBillingProfile(context context.Context, proto *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.UnlinkEmailFromBillingProfile == nil {
		panic("organizationCallbacks.UnlinkEmailFromBillingProfile is not set")
	}
	return organizationCallbacks.UnlinkEmailFromBillingProfile(context, proto)
}

func (MockOrganizationService) LinkLocationToBillingProfile(context context.Context, proto *organizationpb.LinkLocationToBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.LinkLocationToBillingProfile == nil {
		panic("organizationCallbacks.LinkLocationToBillingProfile is not set")
	}
	return organizationCallbacks.LinkLocationToBillingProfile(context, proto)
}

func (MockOrganizationService) UnlinkLocationFromBillingProfile(context context.Context, proto *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error) {
	if organizationCallbacks.UnlinkLocationFromBillingProfile == nil {
		panic("organizationCallbacks.UnlinkLocationFromBillingProfile is not set")
	}
	return organizationCallbacks.UnlinkLocationFromBillingProfile(context, proto)
}
func (MockOrganizationService) RefreshArr(context context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshArr == nil {
		panic("organizationCallbacks.RefreshArr is not set")
	}
	return organizationCallbacks.RefreshArr(context, proto)
}

func (MockOrganizationService) RefreshRenewalSummary(context context.Context, proto *organizationpb.RefreshRenewalSummaryGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshRenewalSummary == nil {
		panic("organizationCallbacks.RefreshRenewalSummary is not set")
	}
	return organizationCallbacks.RefreshRenewalSummary(context, proto)
}
