package events_platform

import (
	"context"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
)

type MockOrganizationServiceCallbacks struct {
	CreateOrganization            func(context.Context, *organizationpb.UpsertOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	AddParent                     func(context.Context, *organizationpb.AddParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	RemoveParent                  func(context.Context, *organizationpb.RemoveParentOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkEmailToOrganization       func(context context.Context, proto *organizationpb.LinkEmailToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
	LinkPhoneNumberToOrganization func(context context.Context, proto *organizationpb.LinkPhoneNumberToOrganizationGrpcRequest) (*organizationpb.OrganizationIdGrpcResponse, error)
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
