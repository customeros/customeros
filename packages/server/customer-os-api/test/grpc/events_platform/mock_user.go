package events_platform

import (
	"context"
	userProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
)

type MockUserServiceCallbacks struct {
	LinkJobRoleToUser     func(context.Context, *userProto.LinkJobRoleToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error)
	LinkPhoneNumberToUser func(context.Context, *userProto.LinkPhoneNumberToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error)
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

func (MockUserService) LinkPhoneNumberToUser(context context.Context, proto *userProto.LinkPhoneNumberToUserGrpcRequest) (*userProto.UserIdGrpcResponse, error) {
	if userCallbacks.LinkPhoneNumberToUser == nil {
		panic("UserCallbacks.CreateUser is not set")
	}
	return userCallbacks.LinkPhoneNumberToUser(context, proto)
}
