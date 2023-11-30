package events_platform

import (
	"context"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
)

type MockContractServiceCallbacks struct {
	CreateContract func(context.Context, *contractpb.CreateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error)
	UpdateContract func(context.Context, *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error)
}

var contractCallbacks = &MockContractServiceCallbacks{}

func SetContractCallbacks(callbacks *MockContractServiceCallbacks) {
	contractCallbacks = callbacks
}

type MockContractService struct {
	contractpb.UnimplementedContractGrpcServiceServer
}

func (MockContractService) CreateContract(context context.Context, proto *contractpb.CreateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	if contractCallbacks.CreateContract == nil {
		panic("contractCallbacks.CreateContract is not set")
	}
	return contractCallbacks.CreateContract(context, proto)
}

func (MockContractService) UpdateContract(context context.Context, proto *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	if contractCallbacks.UpdateContract == nil {
		panic("contractCallbacks.UpdateContract is not set")
	}
	return contractCallbacks.UpdateContract(context, proto)
}
