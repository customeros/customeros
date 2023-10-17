package events_paltform

import (
	"context"
	phonenumberproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
)

type MockPhoneNumberServiceCallbacks struct {
	UpsertPhoneNumber func(ctx context.Context, data *phonenumberproto.UpsertPhoneNumberGrpcRequest) (*phonenumberproto.PhoneNumberIdGrpcResponse, error)
}

var phoneNumberCallbacks = &MockPhoneNumberServiceCallbacks{}

func SetPhoneNumberCallbacks(callbacks *MockPhoneNumberServiceCallbacks) {
	phoneNumberCallbacks = callbacks
}

type MockPhoneNumberService struct {
	phonenumberproto.UnimplementedPhoneNumberGrpcServiceServer
}

func (MockPhoneNumberService) UpsertPhoneNumber(ctx context.Context, data *phonenumberproto.UpsertPhoneNumberGrpcRequest) (*phonenumberproto.PhoneNumberIdGrpcResponse, error) {
	if phoneNumberCallbacks.UpsertPhoneNumber == nil {
		panic("phoneNumberCallbacks.UpsertPhoneNumber is not set")
	}
	return phoneNumberCallbacks.UpsertPhoneNumber(ctx, data)
}
