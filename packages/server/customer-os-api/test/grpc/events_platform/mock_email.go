package events_platform

import (
	"context"
	emailproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
)

type MockEmailServiceCallbacks struct {
	UpsertEmail func(ctx context.Context, data *emailproto.UpsertEmailGrpcRequest) (*emailproto.EmailIdGrpcResponse, error)
}

var emailCallbacks = &MockEmailServiceCallbacks{}

func SetEmailCallbacks(callbacks *MockEmailServiceCallbacks) {
	emailCallbacks = callbacks
}

type MockEmailService struct {
	emailproto.UnimplementedEmailGrpcServiceServer
}

func (MockEmailService) UpsertEmail(ctx context.Context, data *emailproto.UpsertEmailGrpcRequest) (*emailproto.EmailIdGrpcResponse, error) {
	if emailCallbacks.UpsertEmail == nil {
		panic("emailCallbacks.UpsertEmail is not set")
	}
	return emailCallbacks.UpsertEmail(ctx, data)
}
