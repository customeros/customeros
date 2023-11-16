package grpc_client

import (
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Clients struct {
	OrganizationClient organizationpb.OrganizationGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient: organizationpb.NewOrganizationGrpcServiceClient(conn),
	}
	return &clients
}
