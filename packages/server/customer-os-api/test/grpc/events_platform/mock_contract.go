package events_platform

import (
	"context"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
)

type MockContractServiceCallbacks struct {
	RolloutRenewalOpportunityOnExpiration func(context.Context, *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error)
}

var contractCallbacks = &MockContractServiceCallbacks{}

func SetContractCallbacks(callbacks *MockContractServiceCallbacks) {
	contractCallbacks = callbacks
}

type MockContractService struct {
	contractpb.UnimplementedContractGrpcServiceServer
}

func (MockContractService) RolloutRenewalOpportunityOnExpiration(context context.Context, proto *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	if contractCallbacks.RolloutRenewalOpportunityOnExpiration == nil {
		panic("contractCallbacks.RolloutRenewalOpportunityOnExpiration is not set")
	}
	return contractCallbacks.RolloutRenewalOpportunityOnExpiration(context, proto)
}
