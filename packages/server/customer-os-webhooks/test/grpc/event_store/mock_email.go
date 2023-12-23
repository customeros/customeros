package event_store

import (
	"context"
	emailProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
)

type MockEmailServiceCallbacks struct {
	UpsertEmail func(ctx context.Context, data *emailProto.UpsertEmailGrpcRequest) (*emailProto.EmailIdGrpcResponse, error)
}

var emailCallbacks = &MockEmailServiceCallbacks{}

func SetEmailCallbacks(callbacks *MockEmailServiceCallbacks) {
	emailCallbacks = callbacks
}

type MockEmailService struct {
	emailProto.UnimplementedEmailGrpcServiceServer
}

func (MockEmailService) UpsertEmail(ctx context.Context, data *emailProto.UpsertEmailGrpcRequest) (*emailProto.EmailIdGrpcResponse, error) {
	if emailCallbacks.UpsertEmail == nil {
		panic("emailCallbacks.UpsertEmail is not set")
	}
	return emailCallbacks.UpsertEmail(ctx, data)
}
