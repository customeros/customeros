package events_processing_client

import (
	orggrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Client struct {
	OrganizationClient orggrpc.OrganizationGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Client {
	if conn == nil {
		return &Client{}
	}
	clients := Client{
		OrganizationClient: orggrpc.NewOrganizationGrpcServiceClient(conn),
	}
	return &clients
}
