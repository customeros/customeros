package events_processing_client

import (
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Client struct {
	OrganizationClient organizationpb.OrganizationGrpcServiceClient
	ContractClient     contractpb.ContractGrpcServiceClient
	OpportunityClient  opportunitypb.OpportunityGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Client {
	if conn == nil {
		return &Client{}
	}
	clients := Client{
		OrganizationClient: organizationpb.NewOrganizationGrpcServiceClient(conn),
		ContractClient:     contractpb.NewContractGrpcServiceClient(conn),
		OpportunityClient:  opportunitypb.NewOpportunityGrpcServiceClient(conn),
	}
	return &clients
}
