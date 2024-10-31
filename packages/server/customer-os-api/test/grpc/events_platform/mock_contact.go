package events_platform

import (
	"context"
	contactproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
)

type MockContactServiceCallbacks struct {
	LinkPhoneNumberToContact func(context context.Context, proto *contactproto.LinkPhoneNumberToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error)
}

var contactCallbacks = &MockContactServiceCallbacks{}

func SetContactCallbacks(callbacks *MockContactServiceCallbacks) {
	contactCallbacks = callbacks
}

type MockContactService struct {
	contactproto.UnimplementedContactGrpcServiceServer
}

func (MockContactService) LinkPhoneNumberToContact(context context.Context, proto *contactproto.LinkPhoneNumberToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error) {
	if contactCallbacks.LinkPhoneNumberToContact == nil {
		panic("contactCallbacks.LinkPhoneNumberToContact is not set")
	}
	return contactCallbacks.LinkPhoneNumberToContact(context, proto)
}
