package grpc_client

import (
	"google.golang.org/grpc"
)

type Clients struct {
	InvoiceClient invoice_grpc_service.InvoiceGrpcServiceClient
}

func InitClients(conn *grpc.ClientConn) *Clients {
	if conn == nil {
		return &Clients{}
	}
	clients := Clients{
		InvoiceClient: invoice_grpc_service.NewInvoiceGrpcServiceClient(conn),
	}
	return &clients
}
