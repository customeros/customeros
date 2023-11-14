package events_platform

import (
	"context"
	contractproto "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
)

type MockContractServiceCallbacks struct {
	CreateContract func(context.Context, *contractproto.CreateContractGrpcRequest) (*contractproto.ContractIdGrpcResponse, error)
}

var contractCallbacks = &MockContractServiceCallbacks{}

func SetContractCallbacks(callbacks *MockContractServiceCallbacks) {
	contractCallbacks = callbacks
}

type MockContractService struct {
	contractproto.UnimplementedContractGrpcServiceServer
}

func (MockContractService) CreateContract(context context.Context, proto *contractproto.CreateContractGrpcRequest) (*contractproto.ContractIdGrpcResponse, error) {
	if contractCallbacks.CreateContract == nil {
		panic("contractCallbacks.CreateContract is not set")
	}
	return contractCallbacks.CreateContract(context, proto)
}
