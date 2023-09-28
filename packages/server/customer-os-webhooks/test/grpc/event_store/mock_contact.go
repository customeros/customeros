package event_store

import (
	"context"
	contactProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
)

type MockContactServiceCallbacks struct {
	CreateContact      func(context.Context, *contactProto.CreateContactGrpcRequest) (*contactProto.CreateContactGrpcResponse, error)
	LinkEmailToContact func(context context.Context, proto *contactProto.LinkEmailToContactGrpcRequest) (*contactProto.ContactIdGrpcResponse, error)
}

var contactCallbacks = &MockContactServiceCallbacks{}

func SetContactCallbacks(callbacks *MockContactServiceCallbacks) {
	contactCallbacks = callbacks
}

type MockContactService struct {
	contactProto.UnimplementedContactGrpcServiceServer
}

func (MockContactService) CreateContact(context context.Context, proto *contactProto.CreateContactGrpcRequest) (*contactProto.CreateContactGrpcResponse, error) {
	if contactCallbacks.CreateContact == nil {
		panic("contactCallbacks.CreateContact is not set")
	}
	return contactCallbacks.CreateContact(context, proto)
}

func (MockContactService) LinkEmailToContact(context context.Context, proto *contactProto.LinkEmailToContactGrpcRequest) (*contactProto.ContactIdGrpcResponse, error) {
	if contactCallbacks.LinkEmailToContact == nil {
		panic("contactCallbacks.LinkEmailToContact is not set")
	}
	return contactCallbacks.LinkEmailToContact(context, proto)
}
