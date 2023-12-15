package mocked_grpc

import (
	"context"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
)

type MockOrganizationServiceCallbacks struct {
	CreateOrganization            func(ctx context.Context, proto *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	AddParent                     func(ctx context.Context, proto *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RemoveParent                  func(ctx context.Context, proto *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkEmailToOrganization       func(ctx context.Context, proto *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkPhoneNumberToOrganization func(ctx context.Context, proto *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RefreshLastTouchpoint         func(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	UpdateOnboardingStatus        func(ctx context.Context, proto *organizationpb.UpdateOnboardingStatusGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
}

var organizationCallbacks = &MockOrganizationServiceCallbacks{}

func SetOrganizationCallbacks(callbacks *MockOrganizationServiceCallbacks) {
	organizationCallbacks = callbacks
}

type MockOrganizationService struct {
	organizationpb.UnimplementedOrganizationGrpcServiceServer
}

func (MockOrganizationService) UpsertOrganization(ctx context.Context, proto *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.CreateOrganization == nil {
		panic("organizationCallbacks.CreateOrganization is not set")
	}
	return organizationCallbacks.CreateOrganization(ctx, proto)
}

func (MockOrganizationService) LinkEmailToOrganization(ctx context.Context, proto *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.LinkEmailToOrganization == nil {
		panic("organizationCallbacks.LinkEmailToOrganization is not set")
	}
	return organizationCallbacks.LinkEmailToOrganization(ctx, proto)
}

func (MockOrganizationService) LinkPhoneNumberToOrganization(ctx context.Context, proto *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.LinkPhoneNumberToOrganization == nil {
		panic("organizationCallbacks.LinkPhoneNumberToOrganization is not set")
	}
	return organizationCallbacks.LinkPhoneNumberToOrganization(ctx, proto)
}

func (MockOrganizationService) RefreshLastTouchpoint(ctx context.Context, proto *organizationpb.OrganizationIdGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.RefreshLastTouchpoint == nil {
		panic("organizationCallbacks.RefreshLastTouchpoint is not set")
	}
	return organizationCallbacks.RefreshLastTouchpoint(ctx, proto)
}

func (MockOrganizationService) AddParentOrganization(ctx context.Context, proto *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error) {
	if organizationCallbacks.AddParent == nil {
		panic("organizationCallbacks.AddParent is not set")
	}
	return organizationCallbacks.AddParent(ctx, proto)
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
