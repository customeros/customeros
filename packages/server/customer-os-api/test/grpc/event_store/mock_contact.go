package event_store

import (
	"context"
	contactProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
)

type MockContactServiceCallbacks struct {
	CreateContact func(context.Context, *contactProto.CreateContactGrpcRequest) (*contactProto.CreateContactGrpcResponse, error)
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
