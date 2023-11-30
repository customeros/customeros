package events_platform

import (
	"context"
	jobRoleProto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/job_role"
)

type MockJobRoleServiceCallbacks struct {
	CreateJobRole func(context.Context, *jobRoleProto.CreateJobRoleGrpcRequest) (*jobRoleProto.JobRoleIdGrpcResponse, error)
}

var jobRoleCallbacks = &MockJobRoleServiceCallbacks{}

func SetJobRoleCallbacks(callbacks *MockJobRoleServiceCallbacks) {
	jobRoleCallbacks = callbacks
}

type MockJobRoleService struct {
	jobRoleProto.UnimplementedJobRoleGrpcServiceServer
}

func (MockJobRoleService) CreateJobRole(context context.Context, proto *jobRoleProto.CreateJobRoleGrpcRequest) (*jobRoleProto.JobRoleIdGrpcResponse, error) {
	if jobRoleCallbacks.CreateJobRole == nil {
		panic("jobRoleCallbacks.CreateJobRole is not set")
	}
	return jobRoleCallbacks.CreateJobRole(context, proto)
}
