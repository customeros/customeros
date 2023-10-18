package events_paltform

import (
	"context"
	contactproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contact"
)

type MockContactServiceCallbacks struct {
	CreateContact            func(context.Context, *contactproto.UpsertContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error)
	LinkEmailToContact       func(context context.Context, proto *contactproto.LinkEmailToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error)
	LinkPhoneNumberToContact func(context context.Context, proto *contactproto.LinkPhoneNumberToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error)
}

var contactCallbacks = &MockContactServiceCallbacks{}

func SetContactCallbacks(callbacks *MockContactServiceCallbacks) {
	contactCallbacks = callbacks
}

type MockContactService struct {
	contactproto.UnimplementedContactGrpcServiceServer
}

func (MockContactService) UpsertContact(context context.Context, proto *contactproto.UpsertContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error) {
	if contactCallbacks.CreateContact == nil {
		panic("contactCallbacks.CreateContact is not set")
	}
	return contactCallbacks.CreateContact(context, proto)
}

func (MockContactService) LinkEmailToContact(context context.Context, proto *contactproto.LinkEmailToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error) {
	if contactCallbacks.LinkEmailToContact == nil {
		panic("contactCallbacks.LinkEmailToContact is not set")
	}
	return contactCallbacks.LinkEmailToContact(context, proto)
}

func (MockContactService) LinkPhoneNumberToContact(context context.Context, proto *contactproto.LinkPhoneNumberToContactGrpcRequest) (*contactproto.ContactIdGrpcResponse, error) {
	if contactCallbacks.LinkPhoneNumberToContact == nil {
		panic("contactCallbacks.LinkPhoneNumberToContact is not set")
	}
	return contactCallbacks.LinkPhoneNumberToContact(context, proto)
}
