package events_platform

import (
	"context"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockContractServiceCallbacks struct {
	SoftDeleteContract                    func(context.Context, *contractpb.SoftDeleteContractGrpcRequest) (*emptypb.Empty, error)
	RolloutRenewalOpportunityOnExpiration func(context.Context, *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error)
}

var contractCallbacks = &MockContractServiceCallbacks{}

func SetContractCallbacks(callbacks *MockContractServiceCallbacks) {
	contractCallbacks = callbacks
}

type MockContractService struct {
	contractpb.UnimplementedContractGrpcServiceServer
}

func (MockContractService) SoftDeleteContract(context context.Context, proto *contractpb.SoftDeleteContractGrpcRequest) (*emptypb.Empty, error) {
	if contractCallbacks.SoftDeleteContract == nil {
		panic("contractCallbacks.SoftDeleteContract is not set")
	}
	return contractCallbacks.SoftDeleteContract(context, proto)
}

func (MockContractService) RolloutRenewalOpportunityOnExpiration(context context.Context, proto *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	if contractCallbacks.RolloutRenewalOpportunityOnExpiration == nil {
		panic("contractCallbacks.RolloutRenewalOpportunityOnExpiration is not set")
	}
	return contractCallbacks.RolloutRenewalOpportunityOnExpiration(context, proto)
}
