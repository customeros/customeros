package grpc_client

import (
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Clients struct {
	OrganizationClient organization_grpc_service.OrganizationGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient: organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
	}
	return &clients
}
