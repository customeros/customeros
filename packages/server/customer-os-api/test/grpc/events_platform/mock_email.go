package events_platform

import (
	emailproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
)

type MockEmailServiceCallbacks struct {
}

var emailCallbacks = &MockEmailServiceCallbacks{}

func SetEmailCallbacks(callbacks *MockEmailServiceCallbacks) {
	emailCallbacks = callbacks
}

type MockEmailService struct {
	emailproto.UnimplementedEmailGrpcServiceServer
}
