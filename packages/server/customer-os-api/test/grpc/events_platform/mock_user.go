package events_platform

import (
	"context"
	userProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
)

type MockUserServiceCallbacks struct {
	LinkJobRoleToUser func(context.Context, *userProto.LinkJobRoleToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error)
}

var userCallbacks = &MockUserServiceCallbacks{}

func SetUserCallbacks(callbacks *MockUserServiceCallbacks) {
	userCallbacks = callbacks
}

type MockUserService struct {
	userProto.UnimplementedUserGrpcServiceServer
}

func (MockUserService) LinkJobRoleToUser(context context.Context, proto *userProto.LinkJobRoleToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error) {
	if userCallbacks.LinkJobRoleToUser == nil {
		panic("UserCallbacks.CreateUser is not set")
	}
	return userCallbacks.LinkJobRoleToUser(context, proto)
}
