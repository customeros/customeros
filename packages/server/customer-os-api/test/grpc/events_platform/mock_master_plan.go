package events_platform

import (
	"context"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
)

type MockMasterPlanServiceCallbacks struct {
	CreateMasterPlan func(context.Context, *masterplanpb.CreateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error)
}

var masterPlanCallbacks = &MockMasterPlanServiceCallbacks{}

func SetMasterPlanCallbacks(callbacks *MockMasterPlanServiceCallbacks) {
	masterPlanCallbacks = callbacks
}

type MockMasterPlanService struct {
	masterplanpb.UnimplementedMasterPlanGrpcServiceServer
}

func (MockMasterPlanService) CreateMasterPlan(context context.Context, proto *masterplanpb.CreateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error) {
	if masterPlanCallbacks.CreateMasterPlan == nil {
		panic("masterPlanCallbacks.CreateMasterPlan is not set")
	}
	return masterPlanCallbacks.CreateMasterPlan(context, proto)
}
