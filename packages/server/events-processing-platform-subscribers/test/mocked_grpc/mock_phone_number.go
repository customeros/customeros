package mocked_grpc

import (
	"context"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
)

type MockPhoneNumberServiceCallbacks struct {
	RequestPhoneNumberValidation func(ctx context.Context, proto *phonenumberpb.RequestPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error)
}

var PhoneNumberCallbacks = &MockPhoneNumberServiceCallbacks{}

func SetPhoneNumberCallbacks(callbacks *MockPhoneNumberServiceCallbacks) {
	PhoneNumberCallbacks = callbacks
}

type MockPhoneNumberService struct {
	phonenumberpb.UnimplementedPhoneNumberGrpcServiceServer
}

func (MockPhoneNumberService) RequestPhoneNumberValidation(ctx context.Context, proto *phonenumberpb.RequestPhoneNumberValidationGrpcRequest) (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
	if PhoneNumberCallbacks.RequestPhoneNumberValidation == nil {
		panic("PhoneNumberCallbacks.RequestPhoneNumberValidation is not set")
	}
	return PhoneNumberCallbacks.RequestPhoneNumberValidation(ctx, proto)
}
