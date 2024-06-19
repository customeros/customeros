package mocked_grpc

import (
	"context"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
)

type MockEmailServiceCallbacks struct {
	RequestEmailValidation func(ctx context.Context, proto *emailpb.RequestEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error)
}

var EmailCallbacks = &MockEmailServiceCallbacks{}

func SetEmailCallbacks(callbacks *MockEmailServiceCallbacks) {
	EmailCallbacks = callbacks
}

type MockEmailService struct {
	emailpb.UnimplementedEmailGrpcServiceServer
}

func (MockEmailService) RequestEmailValidation(ctx context.Context, proto *emailpb.RequestEmailValidationGrpcRequest) (*emailpb.EmailIdGrpcResponse, error) {
	if EmailCallbacks.RequestEmailValidation == nil {
		panic("EmailCallbacks.RequestEmailValidation is not set")
	}
	return EmailCallbacks.RequestEmailValidation(ctx, proto)
}
