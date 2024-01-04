package events_platform

import (
	"context"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
)

type MockMasterPlanServiceCallbacks struct {
	CreateMasterPlan          func(context.Context, *masterplanpb.CreateMasterPlanGrpcRequest) (*masterplanpb.MasterPlanIdGrpcResponse, error)
	CreateMasterPlanMilestone func(context.Context, *masterplanpb.CreateMasterPlanMilestoneGrpcRequest) (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error)
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

func (MockMasterPlanService) CreateMasterPlanMilestone(context context.Context, proto *masterplanpb.CreateMasterPlanMilestoneGrpcRequest) (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
	if masterPlanCallbacks.CreateMasterPlanMilestone == nil {
		panic("masterPlanCallbacks.CreateMasterPlanMilestone is not set")
	}
	return masterPlanCallbacks.CreateMasterPlanMilestone(context, proto)
}
