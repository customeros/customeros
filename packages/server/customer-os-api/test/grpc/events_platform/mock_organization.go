package events_platform

import (
	"context"

	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
)

type MockOrganizationServiceCallbacks struct {
	CreateOrganization            func(context.Context, *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	AddParent                     func(context.Context, *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RemoveParent                  func(context.Context, *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkEmailToOrganization       func(context context.Context, proto *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkPhoneNumberToOrganization func(context context.Context, proto *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RefreshLastTouchpoint         func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	UpdateOnboardingStatus        func(ctx context.Context, proto *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	UpdateOrganizationOwner       func(ctx context.Context, proto *organizationpb.UpdateOrganizationOwnerGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	CreateBillingProfile          func(ctx context.Context, proto *organizationpb.CreateBillingProfileGrpcRequest) (*organizationpb.BillingProfileIdGrpcResponse, error)
}

var organizationCallbacks = &MockOrganizationServiceCallbacks{}

func SetOrganizationCallbacks(callbacks *MockOrganizationServiceCallbacks) {
	organizationCallbacks = callbacks
}

type MockOrganizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
}

func (MockOrganizationService) UpsertOrganization(context context.Context, proto *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.CreateOrganization == nil {
		panic("organizationCallbacks.CreateOrganization is not set")
	}
	return organizationCallbacks.CreateOrganization(context, proto)
}

func (MockOrganizationService) LinkEmailToOrganization(context context.Context, proto *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.LinkEmailToOrganization == nil {
		panic("organizationCallbacks.LinkEmailToOrganization is not set")
	}
	return organizationCallbacks.LinkEmailToOrganization(context, proto)
}

func (MockOrganizationService) LinkPhoneNumberToOrganization(context context.Context, proto *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.LinkPhoneNumberToOrganization == nil {
		panic("organizationCallbacks.LinkPhoneNumberToOrganization is not set")
	}
	return organizationCallbacks.LinkPhoneNumberToOrganization(context, proto)
}

func (MockOrganizationService) RefreshLastTouchpoint(context context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshLastTouchpoint == nil {
		panic("organizationCallbacks.RefreshLastTouchpoint is not set")
	}
	return organizationCallbacks.RefreshLastTouchpoint(context, proto)
}

func (MockOrganizationService) AddParentOrganization(context context.Context, proto *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.AddParent == nil {
		panic("organizationCallbacks.AddParent is not set")
	}
	return organizationCallbacks.AddParent(context, proto)
}

func (MockOrganizationService) RemoveParentOrganization(context context.Context, proto *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RemoveParent == nil {
		panic("organizationCallbacks.RemoveParent is not set")
	}
	return organizationCallbacks.RemoveParent(context, proto)
}

func (MockOrganizationService) UpdateOnboardingStatus(context context.Context, proto *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.UpdateOnboardingStatus == nil {
		panic("organizationCallbacks.UpdateOnboardingStatus is not set")
	}
	return organizationCallbacks.UpdateOnboardingStatus(context, proto)
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
