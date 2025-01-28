package grpc_client

import (
	invoice_grpc_service "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	organization_grpc_service "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"google.golang.org/grpc"
)

type Clients struct {
	InvoiceClient      invoice_grpc_service.InvoiceGrpcServiceClient
	OrganizationClient organization_grpc_service.OrganizationGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		OrganizationClient: organization_grpc_service.NewOrganizationGrpcServiceClient(conn),
		InvoiceClient:      invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
	}
	return &clients
}
